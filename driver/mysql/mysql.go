package mysql

import (
	"fmt"
	"sort"
	"strings"

	"github.com/things-go/ens"
	"github.com/things-go/ens/driver"
	"golang.org/x/exp/maps"
	"gorm.io/gorm"
)

var _ driver.Driver = (*MySQL)(nil)

type MySQL struct {
	DB                *gorm.DB
	DbName            string
	TableNames        []string
	DisableCommentTag bool
}

func (self *MySQL) GetSchema() (ens.Schemaer, error) {
	mds, err := self.GetEntityMetadata()
	if err != nil {
		return nil, err
	}
	schemas := make([]ens.MixinEntity, 0, len(mds))
	for _, md := range mds {
		schema, err := self.GetEntity(md)
		if err != nil {
			return nil, err
		}
		schemas = append(schemas, schema)
	}
	return &ens.MixinSchema{
		Name:     self.DbName,
		Entities: schemas,
	}, nil
}

func (self *MySQL) GetEntityMetadata() ([]ens.EntityMetadata, error) {
	var tables []Table

	err := self.DB.Table("information_schema.TABLES").
		Scopes(func(db *gorm.DB) *gorm.DB {
			if len(self.TableNames) > 0 {
				db = db.Where("TABLE_NAME IN (?)", self.TableNames)
			}
			return db.Where("TABLE_SCHEMA = ?", self.DbName)
		}).
		Find(&tables).Error
	if err != nil {
		return nil, err
	}
	rows := make([]ens.EntityMetadata, 0, len(tables))
	for _, v := range tables {
		createTableSQL, err := self.GetEntityDefinition(v.Name)
		if err != nil {
			return nil, err
		}
		rows = append(rows, ens.EntityMetadata{
			Name:       v.Name,
			Comment:    v.Comment,
			Definition: createTableSQL,
		})
	}
	return rows, err
}

func (self *MySQL) GetEntity(tb ens.EntityMetadata) (ens.MixinEntity, error) {
	var columns []*Column         // columns
	var indexes []*Index          // indexes
	var foreignKeys []*ForeignKey // foreign key

	// list the database table column.
	err := self.DB.Raw(`
		SELECT 
			* 
		FROM 
			INFORMATION_SCHEMA.COLUMNS 
		WHERE 
			TABLE_SCHEMA = ? 
			AND TABLE_NAME = ? 
		ORDER BY 
			ORDINAL_POSITION`,
		self.DbName, tb.Name).
		Find(&columns).Error
	if err != nil {
		return nil, err
	}

	// index key list
	err = self.DB.Raw("SHOW KEYS FROM `" + tb.Name + "`").Find(&indexes).Error
	if err != nil {
		return nil, err
	}

	// get table column foreign keys
	err = self.DB.Raw(`
		SELECT 
			TABLE_SCHEMA,
			TABLE_NAME,
			COLUMN_NAME,
			REFERENCED_TABLE_SCHEMA,
			REFERENCED_TABLE_NAME,
			REFERENCED_COLUMN_NAME
		FROM 
			INFORMATION_SCHEMA.KEY_COLUMN_USAGE 
		WHERE 
			TABLE_SCHEMA = ? 
			AND REFERENCED_TABLE_NAME IS NOT NULL 
			AND TABLE_NAME = ?`,
		self.DbName, tb.Name).
		Find(&foreignKeys).Error
	if err != nil {
		return nil, err
	}

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
			b.WriteString("KEY")
		} else {
			if strings.EqualFold(keyName, Primary) {
				b.WriteString("PRIMARY KEY")
			} else {
				b.WriteString("UNIQUE KEY")
			}
		}
		if !strings.EqualFold(keyName, Primary) {
			b.WriteString(" ")
			b.WriteString(fmt.Sprintf("`%s`", keyName))
		}

		b.WriteString(" (`")
		b.WriteString(strings.Join(fields, "`, `"))
		b.WriteString("`) USING ")
		b.WriteString(key.IndexType)

		indexer := ens.Index(keyName).Fields(fields...).Definition(b.String())
		indexers = append(indexers, indexer)
	}

	fielders := make([]ens.Fielder, 0, len(columns))
	for _, col := range columns {
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
	return new(ens.EntityBuilder).
		SetMetadata(tb).
		SetFields(fielders...).
		SetIndexes(indexers...), nil
}

func (self *MySQL) GetEntityDefinition(tbName string) (string, error) {
	var ct CreateTable

	err := self.DB.Raw("SHOW CREATE TABLE `" + tbName + "`").Take(&ct).Error
	if err != nil {
		return "", err
	}
	return rAutoIncrement.ReplaceAllString(ct.SQL, " "), err
}
