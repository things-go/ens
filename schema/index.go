package schema

import (
	"fmt"
	"strings"
)

type Index struct {
	Table       string
	KeyName     string
	PrimaryKey  bool
	Unique      bool
	IsComposite bool
	Priority    int
	IndexType   string
	Columns     []string
}

func (self *Index) IntoMysqlString() string {
	b := strings.Builder{}
	b.Grow(32)
	if !self.Unique {
		b.WriteString("KEY")
	} else {
		if self.PrimaryKey {
			b.WriteString("PRIMARY KEY")
		} else {
			b.WriteString("UNIQUE KEY")
		}
	}
	if !self.PrimaryKey {
		b.WriteString(" ")
		b.WriteString(fmt.Sprintf("`%s`", self.KeyName))
	}
	b.WriteString(" (`")
	b.WriteString(strings.Join(self.Columns, "`, `"))
	b.WriteString("`) USING ")
	b.WriteString(self.IndexType)
	return b.String()
}
