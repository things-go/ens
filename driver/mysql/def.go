package mysql

import (
	"fmt"
	"regexp"
	"strings"

	"ariga.io/atlas/sql/mysql"
	"ariga.io/atlas/sql/schema"
	"github.com/things-go/ens"
	"github.com/things-go/ens/internal/sqlx"
	"github.com/things-go/ens/utils"
)

// \b([(]\d+[)])? 匹配0个或1个(\d+)
var typeDictMatchList = []struct {
	Key     string
	NewType func() *ens.GoType
}{
	{`^(bool)`, ens.BoolType},                                 // bool
	{`^(tinyint)\b[(]1[)] unsigned`, ens.BoolType},            // bool
	{`^(tinyint)\b[(]1[)]`, ens.BoolType},                     // bool
	{`^(tinyint)\b([(]\d+[)])? unsigned`, ens.Uint8Type},      // uint8
	{`^(tinyint)\b([(]\d+[)])?`, ens.Int8Type},                // int8
	{`^(smallint)\b([(]\d+[)])? unsigned`, ens.Uint16Type},    // uint16
	{`^(smallint)\b([(]\d+[)])?`, ens.Int16Type},              // int16
	{`^(mediumint)\b([(]\d+[)])? unsigned`, ens.Uint32Type},   // uint32
	{`^(mediumint)\b([(]\d+[)])?`, ens.Int32Type},             // int32
	{`^(int)\b([(]\d+[)])? unsigned`, ens.Uint32Type},         // uint32
	{`^(int)\b([(]\d+[)])?`, ens.Int32Type},                   // int32
	{`^(integer)\b([(]\d+[)])? unsigned`, ens.Uint32Type},     // uint32
	{`^(integer)\b([(]\d+[)])?`, ens.Int32Type},               // int32
	{`^(bigint)\b([(]\d+[)])? unsigned`, ens.Uint64Type},      // uint64
	{`^(bigint)\b([(]\d+[)])?`, ens.Int64Type},                // int64
	{`^(float)\b([(]\d+,\d+[)])? unsigned`, ens.Float32Type},  // float32
	{`^(float)\b([(]\d+,\d+[)])?`, ens.Float32Type},           // float32
	{`^(double)\b([(]\d+,\d+[)])? unsigned`, ens.Float64Type}, // float64
	{`^(double)\b([(]\d+,\d+[)])?`, ens.Float64Type},          // float64
	{`^(char)\b[(]\d+[)]`, ens.StringType},                    // string
	{`^(varchar)\b[(]\d+[)]`, ens.StringType},                 // string
	{`^(datetime)\b([(]\d+[)])?`, ens.TimeType},               // time.Time
	{`^(date)\b([(]\d+[)])?`, ens.TimeType},                   // datatypes.Date
	{`^(timestamp)\b([(]\d+[)])?`, ens.TimeType},              // time.Time
	{`^(time)\b([(]\d+[)])?`, ens.TimeType},                   // time.Time
	{`^(year)\b([(]\d+[)])?`, ens.TimeType},                   // time.Time
	{`^(text)\b([(]\d+[)])?`, ens.StringType},                 // string
	{`^(tinytext)\b([(]\d+[)])?`, ens.StringType},             // string
	{`^(mediumtext)\b([(]\d+[)])?`, ens.StringType},           // string
	{`^(longtext)\b([(]\d+[)])?`, ens.StringType},             // string
	{`^(blob)\b([(]\d+[)])?`, ens.BytesType},                  // []byte
	{`^(tinyblob)\b([(]\d+[)])?`, ens.BytesType},              // []byte
	{`^(mediumblob)\b([(]\d+[)])?`, ens.BytesType},            // []byte
	{`^(longblob)\b([(]\d+[)])?`, ens.BytesType},              // []byte
	{`^(bit)\b[(]\d+[)]`, ens.BytesType},                      // []uint8
	{`^(json)\b`, ens.JSONRawMessageType},                     // datatypes.JSON
	{`^(enum)\b[(](.)+[)]`, ens.StringType},                   // string
	{`^(set)\b[(](.)+[)]`, ens.StringType},                    // string
	{`^(decimal)\b[(]\d+,\d+[)]`, ens.DecimalType},            // string
	{`^(binary)\b[(]\d+[)]`, ens.BytesType},                   // []byte
	{`^(varbinary)\b[(]\d+[)]`, ens.BytesType},                // []byte
	{`^(geometry)`, ens.StringType},                           // string
}

func intoGoType(columnType string) *ens.GoType {
	for _, v := range typeDictMatchList {
		ok, _ := regexp.MatchString(v.Key, columnType)
		if ok {
			return v.NewType()
		}
	}
	panic(fmt.Sprintf("type (%v) not match in any way, need to add on (https://github.com/things-go/ormat/blob/main/driver/mysql/def.go)", columnType))
}

type TableDef struct {
	tb *schema.Table
}

func NewTableDef(tb *schema.Table) ens.TableDef {
	return &TableDef{tb: tb}
}

func (self *TableDef) Table() *schema.Table { return self.tb }

func (self *TableDef) PrimaryKey() ens.IndexDef {
	if self.tb.PrimaryKey != nil {
		return NewIndexDef(self.tb.PrimaryKey)
	}
	return nil
}

func (self *TableDef) Definition() string {
	tb := self.tb

	b := &strings.Builder{}
	b.Grow(64)
	fmt.Fprintf(b, "CREATE TABLE `%s` (\n", tb.Name)

	remain := len(tb.Columns) + len(tb.Indexes) + len(tb.ForeignKeys)
	if tb.PrimaryKey != nil {
		remain++
	}
	suffixOrEmpty := func(r int) string {
		if r == 0 {
			return ""
		}
		return ","
	}
	//* columns
	for _, col := range tb.Columns {
		remain--
		suffix := suffixOrEmpty(remain)
		comment, ok := sqlx.Comment(col.Attrs)
		if ok {
			comment = fmt.Sprintf(" COMMENT '%s'", comment)
		}
		fmt.Fprintf(b, "  `%s` %s%s%s\n", col.Name, NewColumnDef(col).Definition(), comment, suffix)
	}
	//* pk + indexes
	if tb.PrimaryKey != nil {
		remain--
		suffix := suffixOrEmpty(remain)
		fmt.Fprintf(b, "  %s%s\n", NewIndexDef(tb.PrimaryKey).Definition(), suffix)
	}
	for _, val := range tb.Indexes {
		remain--
		if sqlx.IndexEqual(tb.PrimaryKey, val) { // ignore primary key, maybe include
			continue
		}
		suffix := suffixOrEmpty(remain)
		fmt.Fprintf(b, "  %s%s\n", NewIndexDef(val).Definition(), suffix)
	}
	//* foreignKeys
	for _, val := range tb.ForeignKeys {
		remain--
		suffix := suffixOrEmpty(remain)
		fmt.Fprintf(b, "  %s%s\n", NewForeignKey(val).Definition(), suffix)
	}

	engine := mysql.EngineInnoDB
	charset := "utf8mb4"
	collate := ""
	comment := ""
	for _, attr := range tb.Attrs {
		switch val := attr.(type) {
		case *mysql.Engine:
			engine = val.V
		case *schema.Charset:
			charset = val.V
		case *schema.Collation:
			collate = val.V
		case *schema.Comment:
			comment = val.Text
			// case *mysql.AutoIncrement: // ignore this
		}
	}
	fmt.Fprintf(b, ") ENGINE=%s DEFAULT CHARSET=%s", engine, charset)
	if collate != "" {
		fmt.Fprintf(b, " COLLATE='%s'", collate)
	}
	if comment != "" {
		fmt.Fprintf(b, " COMMENT='%s'", comment)
	}
	return b.String()
}

type ColumnDef struct {
	col *schema.Column
}

func NewColumnDef(col *schema.Column) ens.ColumnDef {
	return &ColumnDef{col: col}
}

func (self *ColumnDef) Column() *schema.Column { return self.col }

func (self *ColumnDef) Definition() string {
	nullable := self.col.Type.Null
	autoIncrement := autoIncrement(self.col.Attrs)

	b := &strings.Builder{}
	b.Grow(64)
	b.WriteString(self.col.Type.Raw)
	if !nullable {
		b.WriteString(" NOT NULL")
	}
	if autoIncrement {
		b.WriteString(" AUTO_INCREMENT")
	} else {
		dv, ok := sqlx.DefaultValue(self.col)
		if ok {
			fmt.Fprintf(b, " DEFAULT '%s'", strings.Trim(dv, `"`))
		} else if nullable {
			b.WriteString(" DEFAULT NULL")
		}
	}
	return b.String()
}

// column, type, not null, authIncrement, default, [primaryKey|index], comment
func (self *ColumnDef) GormTag(tb *schema.Table) string {
	col := self.col

	pkPriority, isPk := 0, false
	if pk := tb.PrimaryKey; pk != nil {
		pkPriority, isPk = sqlx.FindIndexPartSeq(pk.Parts, col)
	}
	autoIncrement := autoIncrement(col.Attrs)

	b := &strings.Builder{}
	b.Grow(64)
	fmt.Fprintf(b, `gorm:"column:%s`, col.Name)
	if !(isPk && autoIncrement) {
		fmt.Fprintf(b, ";type:%s", col.Type.Raw)
	}
	if !col.Type.Null {
		fmt.Fprintf(b, ";not null")
	}

	if isPk {
		if autoIncrement {
			fmt.Fprintf(b, ";autoIncrement:true")
		}
	} else {
		dv, ok := sqlx.DefaultValue(col)
		if ok {
			if dv == `""` || dv == "" {
				dv = "''"
			} else {
				dv = strings.Trim(dv, `"`) // format: `"xxx"` or `'xxx'`
			}
			fmt.Fprintf(b, ";default:%s", dv)
		} else if col.Type.Null {
			fmt.Fprintf(b, ";default:null")
		}
	}

	//* pk + indexes
	if isPk && tb.PrimaryKey != nil {
		fmt.Fprintf(b, ";primaryKey")
		if len(tb.PrimaryKey.Parts) > 1 {
			fmt.Fprintf(b, ",priority:%d", pkPriority)
		}
	}
	for _, val := range col.Indexes {
		if sqlx.IndexEqual(tb.PrimaryKey, val) { // ignore primary key, may be include
			continue
		}
		if val.Unique {
			fmt.Fprintf(b, ";uniqueIndex:%s", val.Name)
		} else {
			fmt.Fprintf(b, ";index:%s", val.Name)
			// 	mysql.IndexTypeFullText
			// if v.IndexType == "FULLTEXT" {
			// 	b.WriteString(",class:FULLTEXT")
			// }
		}
		if len(val.Parts) > 1 {
			priority, ok := sqlx.FindIndexPartSeq(val.Parts, col)
			if ok {
				fmt.Fprintf(b, ",priority:%d", priority)
			}
		}
	}
	if comment, ok := sqlx.Comment(col.Attrs); ok && comment != "" {
		fmt.Fprintf(b, ";comment:%s", utils.TrimFieldComment(comment))
	}
	b.WriteString(`"`)
	return b.String()
}

type IndexDef struct {
	index *schema.Index
}

func NewIndexDef(index *schema.Index) ens.IndexDef {
	return &IndexDef{index: index}
}

func (self *IndexDef) Index() *schema.Index { return self.index }

func (self *IndexDef) Definition() string {
	index := self.index
	fields := sqlx.IndexPartColumnNames(index.Parts)
	indexType := findIndexType(index.Attrs)
	fieldList := "`" + strings.Join(fields, "`,`") + "`"
	if sqlx.IndexEqual(index.Table.PrimaryKey, index) {
		return fmt.Sprintf("PRIMARY KEY (%s) USING %s", fieldList, indexType)
	} else if index.Unique {
		return fmt.Sprintf("UNIQUE KEY `%s` (%s) USING %s", index.Name, fieldList, indexType)
	} else {
		return fmt.Sprintf("KEY `%s` (%s) USING %s", index.Name, fieldList, indexType)
	}
}

type ForeignKeyDef struct {
	fk *schema.ForeignKey
}

func NewForeignKey(fk *schema.ForeignKey) ens.ForeignKeyDef {
	return &ForeignKeyDef{fk: fk}
}

func (self *ForeignKeyDef) ForeignKey() *schema.ForeignKey { return self.fk }

func (self *ForeignKeyDef) Definition() string {
	fk := self.fk
	columnNameList := "`" + strings.Join(sqlx.ColumnNames(fk.Columns), "`,`") + "`"
	refColumnNameList := "`" + strings.Join(sqlx.ColumnNames(fk.RefColumns), "`,`") + "`"
	return fmt.Sprintf(
		"CONSTRAINT `%s` FOREIGN KEY (%s) REFERENCES `%s` (%s) ON DELETE %s ON UPDATE %s",
		fk.Symbol, columnNameList, fk.RefTable.Name, refColumnNameList, fk.OnDelete, fk.OnUpdate,
	)
}
