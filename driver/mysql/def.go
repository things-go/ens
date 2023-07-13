package mysql

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/things-go/ens"
	"github.com/things-go/ens/utils"
)

const (
	Primary = "PRIMARY"

	nullableTrue  = "YES"
	nullableFalse = "NO"

	extraAutoIncrement = "auto_increment"

	columnKeyPrimary  = "PRI"
	columnKeyUnique   = "UNI"
	columnKeyMultiple = "MUL"
)

type ColumnKeyType int

const (
	ColumnKeyType_NotKey   ColumnKeyType = iota // no key
	ColumnKeyType_Primary                       // primary key
	ColumnKeyType_Multiple                      // multiple index key
	ColumnKeyType_Unique                        // unique
)

// Table mysql table info
// sql: SELECT * FROM information_schema.TABLES WHERE TABLE_SCHEMA={db_name}
type Table struct {
	Name      string `gorm:"column:TABLE_NAME"`    // table name, 表名
	Engine    string `gorm:"column:ENGINE"`        // table engine, 表引擎(InnoDB)
	RowFormat string `gorm:"column:ROW_FORMAT"`    // table row format, 表数据格式(Dynamic)
	Comment   string `gorm:"column:TABLE_COMMENT"` // table comment, 表注释
}

// Column mysql column info
// sql: SELECT * FROM `INFORMATION_SCHEMA`.`COLUMNS` WHERE `TABLE_SCHEMA`={dbName} AND `TABLE_NAME`={tbName}
type Column struct {
	ColumnName             string  `gorm:"column:COLUMN_NAME"`      // column name
	OrdinalPosition        int     `gorm:"column:ORDINAL_POSITION"` // column ordinal position
	ColumnDefault          *string `gorm:"column:COLUMN_DEFAULT"`   // column default value.null mean not set.
	IsNullable             string  `gorm:"column:IS_NULLABLE"`      // column null or not, YEW/NO
	DataType               string  `gorm:"column:DATA_TYPE"`        // column data type(varchar)
	CharacterMaximumLength int64   `gorm:"column:CHARACTER_MAXIMUM_LENGTH"`
	CharacterOctetLength   int64   `gorm:"column:CHARACTER_OCTET_LENGTH"`
	NumericPrecision       int64   `gorm:"column:NUMERIC_PRECISION"`
	NumericScale           int64   `gorm:"column:NUMERIC_SCALE"`
	ColumnType             string  `gorm:"column:COLUMN_TYPE"`    // column type(varchar(64))
	ColumnKey              string  `gorm:"column:COLUMN_KEY"`     // column key, PRI/MUL
	Extra                  string  `gorm:"column:EXTRA"`          // extra (auto_increment)
	ColumnComment          string  `gorm:"column:COLUMN_COMMENT"` // column comment
}

func (c *Column) IntoSqlDefinition() string {
	nullable := strings.EqualFold(c.IsNullable, "YES")
	isAutoIncrement := c.Extra == extraAutoIncrement
	b := strings.Builder{}
	b.Grow(64)
	b.WriteString(c.ColumnType)
	if !nullable {
		b.WriteString(" ")
		b.WriteString("NOT NULL")
	}
	if isAutoIncrement {
		b.WriteString(" ")
		b.WriteString("AUTO_INCREMENT")
	} else {
		dv := ""
		if c.ColumnDefault != nil {
			dv = fmt.Sprintf("DEFAULT '%s'", *c.ColumnDefault)
		} else if nullable {
			dv = "DEFAULT NULL"
		}
		if dv != "" {
			b.WriteString(" ")
			b.WriteString(dv)
		}
	}
	return b.String()
}

func (c *Column) IntoOrmTag(indexes []*Index, keyNameCount map[string]int, disableCommentTag bool) string {
	nullable := strings.EqualFold(c.IsNullable, "YES")
	isAutoIncrement := c.Extra == extraAutoIncrement

	b := strings.Builder{}
	b.Grow(64)
	b.WriteString(`gorm:"column:`)
	b.WriteString(c.ColumnName)
	// FIXME: 主要是整型主键,gorm在自动迁移时没有在mysql上加上auto_increment
	if !(c.ColumnKey == "PRI" && isAutoIncrement) {
		b.WriteString(";")
		b.WriteString("type:")
		b.WriteString(c.ColumnType)
	}
	if !nullable {
		b.WriteString(";")
		b.WriteString("not null")
	}
	if isAutoIncrement {
		b.WriteString(";")
		b.WriteString("autoIncrement:true")
	} else {
		if c.ColumnDefault != nil {
			b.WriteString(";")
			if *c.ColumnDefault == "" {
				b.WriteString("default:''")
			} else {
				b.WriteString("default:")
				b.WriteString(*c.ColumnDefault)
			}
		} else if nullable {
			b.WriteString(";")
			b.WriteString("default:null")
		}
	}

	for _, v := range indexes {
		isPrimaryKey := false
		b.WriteString(";")
		if v.NonUnique {
			if v.KeyName == "sort" { // 兼容 gorm 本身 sort 标签
				b.WriteString("index")
			} else {
				b.WriteString("index:")
				b.WriteString(v.KeyName)
			}
			if v.IndexType == "FULLTEXT" {
				b.WriteString(",class:FULLTEXT")
			}
		} else {
			if strings.EqualFold(v.KeyName, Primary) {
				b.WriteString("primaryKey")
				isPrimaryKey = true
			} else {
				b.WriteString("uniqueIndex:")
				b.WriteString(v.KeyName)
			}
		}
		if keyNameCount[v.KeyName] > 1 {
			if isPrimaryKey {
				b.WriteString(";")
			} else {
				b.WriteString(",")
			}
			b.WriteString("priority:")
			b.WriteString(strconv.FormatInt(int64(v.SeqInIndex), 10))
		}
	}

	if c.ColumnComment != "" && !disableCommentTag {
		b.WriteString(";")
		b.WriteString("comment:")
		b.WriteString(utils.TrimFieldComment(c.ColumnComment))
	}
	b.WriteString(`"`)
	return b.String()
}

// key index info
// sql: SHOW KEYS FROM {table_name}
type Index struct {
	Table      string `gorm:"column:Table"`        // 表名
	NonUnique  bool   `gorm:"column:Non_unique"`   // 不是唯一索引
	KeyName    string `gorm:"column:Key_name"`     // 索引关键字
	SeqInIndex int    `gorm:"column:Seq_in_index"` // 索引排序
	ColumnName string `gorm:"column:Column_name"`  // 索引列名
	IndexType  string `gorm:"column:Index_type"`   // 索引类型, BTREE
}

type IndexSlice []*Index

func (t IndexSlice) Len() int           { return len(t) }
func (t IndexSlice) Less(i, j int) bool { return t[i].SeqInIndex < t[j].SeqInIndex }
func (t IndexSlice) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }

// ForeignKey Foreign key of db table info . 表的外键信息
// sql: SELECT table_schema, table_name, column_name, referenced_table_schema, referenced_table_name, referenced_column_name
//
//	FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE
//	WHERE table_schema={db_name} AND REFERENCED_TABLE_NAME IS NOT NULL AND TABLE_NAME={table_name}
type ForeignKey struct {
	TableSchema           string `gorm:"column:table_schema"`            // Database of column.
	TableName             string `gorm:"column:table_name"`              // Data table of column.
	ColumnName            string `gorm:"column:column_name"`             // column names.
	ReferencedTableSchema string `gorm:"column:referenced_table_schema"` // The database where the index is located.
	ReferencedTableName   string `gorm:"column:referenced_table_name"`   // Affected tables .
	ReferencedColumnName  string `gorm:"column:referenced_column_name"`  // Which column of the affected table.
}

// CreateTable mysql show create table information.
// sql: SHOW CREATE TABLE {tableName}
type CreateTable struct {
	Table string `gorm:"column:Table"`
	SQL   string `gorm:"column:Create Table"`
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

type dictMatchKv struct {
	Key     string
	NewType func() *ens.GoType
}

// \b([(]\d+[)])? 匹配0个或1个(\d+)
var typeDictMatchList = []dictMatchKv{
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
	{`^(text)\b([(]\d+[)])?`, ens.StringType},                 // string
	{`^(tinytext)\b([(]\d+[)])?`, ens.StringType},             // string
	{`^(mediumtext)\b([(]\d+[)])?`, ens.StringType},           // string
	{`^(longtext)\b([(]\d+[)])?`, ens.StringType},             // string
	{`^(blob)\b([(]\d+[)])?`, ens.BytesType},                  // []byte
	{`^(tinyblob)\b([(]\d+[)])?`, ens.BytesType},              // []byte
	{`^(mediumblob)\b([(]\d+[)])?`, ens.BytesType},            // []byte
	{`^(longblob)\b([(]\d+[)])?`, ens.BytesType},              // []byte
	{`^(bit)\b[(]\d+[)]`, ens.BytesType},                      // []uint8
	{`^(json)\b`, ens.BytesType},                              // datatypes.JSON
	{`^(enum)\b[(](.)+[)]`, ens.StringType},                   // string
	{`^(decimal)\b[(]\d+,\d+[)]`, ens.DecimalType},            // string
	{`^(binary)\b[(]\d+[)]`, ens.BytesType},                   // []byte
	{`^(varbinary)\b[(]\d+[)]`, ens.BytesType},                // []byte
	{`geometry`, ens.StringType},                              // string
}
