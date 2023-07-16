package schema

import (
	"strconv"
	"strings"

	"github.com/things-go/ens/utils"
)

type Column struct {
	Name              string
	DataType          string // varchar
	ColumnType        string // varchar(255)
	Unique            bool
	Nullable          bool
	Default           *string
	IsPrimaryKey      bool
	AutoIncrement     bool
	HasLength         bool
	Length            int64
	HasPrecisionScale bool
	Precision         int64
	Scale             int64
	Comment           string
	Indexes           []*Index
}

// column, type, not null, authIncrement, default, [primaryKey|index], comment
func (c *Column) IntoOrmTag() string {
	b := strings.Builder{}
	b.Grow(64)
	b.WriteString(`gorm:"column:`)
	b.WriteString(c.Name)

	// FIXME: 主要是整型主键,gorm在自动迁移时没有在mysql上加上auto_increment
	if !(c.IsPrimaryKey && c.AutoIncrement) {
		b.WriteString(";")
		b.WriteString("type:")
		b.WriteString(c.ColumnType)
	}
	if !c.Nullable {
		b.WriteString(";")
		b.WriteString("not null")
	}
	if c.IsPrimaryKey {
		if c.AutoIncrement {
			b.WriteString(";")
			b.WriteString("autoIncrement:true")
		}
	} else {
		if c.Default != nil {
			b.WriteString(";")
			if *c.Default == "" {
				b.WriteString("default:''")
			} else {
				b.WriteString("default:")
				b.WriteString(*c.Default)
			}
		} else if c.Nullable {
			b.WriteString(";")
			b.WriteString("default:null")
		}
	}

	for _, v := range c.Indexes {
		b.WriteString(";")
		if v.PrimaryKey {
			b.WriteString("primaryKey")
		} else if v.Unique {
			b.WriteString("uniqueIndex:")
			b.WriteString(v.KeyName)
		} else {
			if v.KeyName == "sort" { // 兼容 gorm 本身 sort 标签
				b.WriteString("index")
			} else {
				b.WriteString("index:")
				b.WriteString(v.KeyName)
			}
			if v.IndexType == "FULLTEXT" {
				b.WriteString(",class:FULLTEXT")
			}
		}
		if v.IsComposite {
			b.WriteString(",")
			b.WriteString("priority:")
			b.WriteString(strconv.FormatInt(int64(v.Priority), 10))
		}
	}

	if c.Comment != "" {
		b.WriteString(";")
		b.WriteString("comment:")
		b.WriteString(utils.TrimFieldComment(c.Comment))
	}
	b.WriteString(`"`)
	return b.String()
}
