package mysql

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/things-go/ens"
	"github.com/things-go/ens/driver"
	"github.com/xwb1989/sqlparser"
	"golang.org/x/exp/maps"
)

var _ driver.Driver = (*SQL)(nil)

type SQL struct {
	CreateTableSQL    string
	DisableCommentTag bool
	entity            ens.MixinEntity
}

func (self *SQL) hasParse() bool {
	return self.entity != nil
}

func (self *SQL) Parse() error {
	if self.hasParse() {
		return nil
	}
	statement, err := sqlparser.Parse(self.CreateTableSQL)
	if err != nil {
		return err
	}
	switch stmt := statement.(type) {
	case *sqlparser.DDL:
		if stmt.Action != sqlparser.CreateStr {
			return errors.New("不是创建表语句")
		}
		if stmt.TableSpec == nil {
			return errors.New("未解析到任何字段")
		}
		tbName := stmt.NewName.Name.String()
		md := ens.EntityMetadata{
			Name:       tbName,
			Comment:    "",
			Definition: self.CreateTableSQL,
		}
		// ENGINE=InnoDB default charset=utf8mb4 collate=utf8mb4_general_ci comment='我是注释'
		tbOptions := strings.Split(stmt.TableSpec.Options, " ")
		for _, option := range tbOptions {
			keyValue := strings.Split(option, "=")
			if len(keyValue) >= 2 {
				switch keyValue[0] {
				case "ENGINE":
				case "charset":
				case "collate":
				case "comment":
					md.Comment = strings.ReplaceAll(keyValue[1], "'", "")
				}
			}
		}
		//* parser index definition
		indexes := fromSqlIndexDefinition(tbName, stmt.TableSpec.Indexes)

		keyNameCount := make(map[string]int)          // key name count
		columnNameMapKey := make(map[string][]*Index) // column name map key
		keyNameMapFields := make(map[string][]*Index)
		indexers := make([]ens.Indexer, 0, len(indexes))
		for _, key := range indexes {
			keyNameCount[key.KeyName]++
			keyNameMapFields[key.KeyName] = append(keyNameMapFields[key.KeyName], key)
			columnNameMapKey[key.ColumnName] = append(columnNameMapKey[key.ColumnName], key)
		}
		keyNames := maps.Keys(keyNameCount)
		sort.Strings(keyNames)
		for _, keyName := range keyNames {
			keys, ok := keyNameMapFields[keyName]
			if !ok || len(keys) == 0 {
				continue
			}
			sort.Sort(IndexSlice(keys))
			fields := make([]string, 0, len(keys))
			for _, v := range keys {
				fields = append(fields, v.ColumnName)
			}
			b := strings.Builder{}
			b.Grow(32)
			key := keys[0]
			if key.NonUnique {
				b.WriteString("KEY ")
			} else {
				if strings.EqualFold(keyName, Primary) {
					b.WriteString("PRIMARY KEY ")
				} else {
					b.WriteString("UNIQUE KEY ")
				}
			}
			b.WriteString(fmt.Sprintf("`%s`", keyName))
			b.WriteString(" (`")
			b.WriteString(strings.Join(fields, "`, `"))
			b.WriteString("`) USING ")
			b.WriteString(key.IndexType)

			indexer := ens.Index(keyName).Fields(fields...).Definition(b.String())
			indexers = append(indexers, indexer)
		}

		fielders := make([]ens.Fielder, 0, len(stmt.TableSpec.Columns))
		for i := 0; i < len(stmt.TableSpec.Columns); i++ {
			//* parser column definition
			col, err := fromSqlColumnDefinition(i+1, stmt.TableSpec.Columns[i])
			if err != nil {
				return err
			}
			col.ColumnKey = findColumnKey(indexes, col.ColumnName)

			nullable := strings.EqualFold(col.IsNullable, nullableTrue)
			fielder := ens.
				Field(intoGoType(col.ColumnType), col.ColumnName).
				Comment(col.ColumnComment).
				Tags(col.IntoOrmTag(columnNameMapKey[col.ColumnName], keyNameCount, self.DisableCommentTag)).
				Definition(col.IntoSqlDefinition())
			if nullable {
				fielder.Nullable().Optional()
			}
			fielders = append(fielders, fielder)
		}

		self.entity = new(ens.EntityBuilder).
			SetMetadata(md).
			SetFields(fielders...).
			SetIndexes(indexers...)
	default:
		return errors.New("不是DDL语句")
	}
	return nil
}

func (self *SQL) GetSchema() (ens.Schemaer, error) {
	err := self.Parse()
	if err != nil {
		return nil, err
	}
	return &ens.MixinSchema{
		Name:     "",
		Entities: []ens.MixinEntity{self.entity},
	}, nil
}

func (self *SQL) GetEntityMetadata() ([]ens.EntityMetadata, error) {
	err := self.Parse()
	if err != nil {
		return nil, err
	}
	return []ens.EntityMetadata{self.entity.Metadata()}, nil
}

func (self *SQL) GetEntity(tb ens.EntityMetadata) (ens.MixinEntity, error) {
	err := self.Parse()
	if err != nil {
		return nil, err
	}
	return self.entity, nil
}

func (self *SQL) GetEntityDefinition(tbName string) (string, error) {
	return self.CreateTableSQL, nil
}

func fromSqlIndexDefinition(tbName string, idxs []*sqlparser.IndexDefinition) []*Index {
	indexes := make([]*Index, 0, 8)
	for _, idx := range idxs {
		keyName := idx.Info.Name.String()
		nonUnique := true
		switch idx.Info.Type {
		case "primary key":
			nonUnique = false
		case "unique key", "unique":
			nonUnique = false
		case "key":
			nonUnique = true
		}
		indexType := "BTREE"
		for _, option := range idx.Options {
			if option.Name == "using" {
				indexType = option.Using
			}
		}

		for i, col := range idx.Columns {
			indexes = append(indexes, &Index{
				Table:      tbName,
				NonUnique:  nonUnique,
				KeyName:    keyName,
				SeqInIndex: i + 1,
				ColumnName: col.Column.String(),
				IndexType:  indexType,
			})
		}
	}
	return indexes
}

func findColumnKey(indexes []*Index, columnName string) string {
	ck := ""
	for _, v := range indexes {
		if v.ColumnName == columnName {
			if v.KeyName == Primary {
				return columnKeyPrimary
			}
			if !v.NonUnique {
				return columnKeyUnique
			} else {
				ck = columnKeyMultiple
			}
		}
	}
	return ck
}

func fromSqlColumnDefinition(ordinalPosition int, col *sqlparser.ColumnDefinition) (*Column, error) {
	colType := col.Type
	column := &Column{
		ColumnName:             col.Name.String(),
		OrdinalPosition:        ordinalPosition,
		ColumnDefault:          nil,
		IsNullable:             "",
		DataType:               colType.Type,
		CharacterMaximumLength: 0,
		CharacterOctetLength:   0,
		NumericPrecision:       0,
		NumericScale:           0,
		ColumnType:             "",
		ColumnKey:              "",
		Extra:                  "",
		ColumnComment:          "",
	}

	if colType.Default != nil {
		defaultValue := string(colType.Default.Val)
		column.ColumnDefault = &defaultValue
	}
	if colType.NotNull {
		column.IsNullable = nullableFalse
	} else {
		column.IsNullable = nullableTrue
	}
	if colType.Autoincrement {
		column.Extra = extraAutoIncrement
	}
	if colType.Comment != nil {
		column.ColumnComment = string(colType.Comment.Val)
	}
	toInt := func(l *sqlparser.SQLVal) int64 {
		if l == nil {
			return 0
		}
		length, err := strconv.ParseInt(string(l.Val), 0, 0)
		if err != nil {
			return 0
		}
		return length
	}

	isUnsigned := bool(colType.Unsigned)
	switch colType.Type {
	case "tinyint", "smallint", "mediumint",
		"int", "integer", "bigint":
		// `^(tinyint)\b[(]1[)] unsigned`
		// `^(tinyint)\b[(]1[)]`
		// `^(tinyint)\b([(]\d+[)])? unsigned`
		// `^(tinyint)\b([(]\d+[)])?`
		// `^(smallint)\b([(]\d+[)])? unsigned`
		// `^(smallint)\b([(]\d+[)])?`
		// `^(mediumint)\b([(]\d+[)])? unsigned`
		// `^(mediumint)\b([(]\d+[)])?`
		// `^(int)\b([(]\d+[)])? unsigned`
		// `^(int)\b([(]\d+[)])?`
		// `^(integer)\b([(]\d+[)])? unsigned`
		// `^(integer)\b([(]\d+[)])?`
		// `^(bigint)\b([(]\d+[)])? unsigned`
		// `^(bigint)\b([(]\d+[)])?`
		length := toInt(colType.Length)
		if length > 0 {
			if isUnsigned {
				column.ColumnType = fmt.Sprintf("%s(%d) unsigned", colType.Type, length)
			} else {
				column.ColumnType = fmt.Sprintf("%s(%d)", colType.Type, length)
			}
		} else {
			if isUnsigned {
				column.ColumnType = fmt.Sprintf("%s unsigned", colType.Type)
			} else {
				column.ColumnType = colType.Type
			}
		}
	case "float", "double", "decimal":
		// `^(float)\b([(]\d+,\d+[)])? unsigned`
		// `^(float)\b([(]\d+,\d+[)])?`
		// `^(double)\b([(]\d+,\d+[)])? unsigned`
		// `^(double)\b([(]\d+,\d+[)])?`
		// `^(decimal)\b[(]\d+,\d+[)]`
		length, scale := toInt(colType.Length), toInt(colType.Scale)
		if length > 0 {
			column.ColumnType = fmt.Sprintf("%s(%d,%d)", colType.Type, length, scale)
		} else {
			column.ColumnType = colType.Type
		}
	case "char", "varchar",
		"text", "tinytext", "mediumtext", "longtext",
		"date", "datetime", "timestamp", "time",
		"blob", "tinyblob", "mediumblob", "longblob",
		"binary", "varbinary",
		"bit":
		// `^(char)\b[(]\d+[)]`
		// `^(varchar)\b[(]\d+[)]`
		// `^(text)\b([(]\d+[)])?`
		// `^(tinytext)\b([(]\d+[)])?`
		// `^(mediumtext)\b([(]\d+[)])?`
		// `^(longtext)\b([(]\d+[)])?`
		// `^(date)\b([(]\d+[)])?`
		// `^(datetime)\b([(]\d+[)])?`
		// `^(timestamp)\b([(]\d+[)])?`
		// `^(time)\b([(]\d+[)])?`
		// `^(blob)\b([(]\d+[)])?`
		// `^(tinyblob)\b([(]\d+[)])?`
		// `^(mediumblob)\b([(]\d+[)])?`
		// `^(longblob)\b([(]\d+[)])?`
		// `^(binary)\b[(]\d+[)]`
		// `^(varbinary)\b[(]\d+[)]`
		// `^(bit)\b[(]\d+[)]`
		length := toInt(colType.Length)
		if length > 0 {
			column.ColumnType = fmt.Sprintf("%s(%d)", colType.Type, length)
		} else {
			column.ColumnType = colType.Type
		}
	case "enum":
		// `^(enum)\b[(](.)+[)]`
		column.ColumnType = "enum(" + strings.Join(colType.EnumValues, ",") + ")"
	case "json":
		// `^(json)\b`
		column.ColumnType = colType.Type
	case "geometry":
		// `geometry`
		column.ColumnType = colType.Type
	default:
		return nil, fmt.Errorf("unsupported column type(%s)", colType.Type)
	}
	return column, nil
}
