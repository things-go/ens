package mysql

import (
	"ariga.io/atlas/sql/mysql"
	"ariga.io/atlas/sql/schema"

	"github.com/things-go/ens/internal/insql"
)

func autoIncrement(attrs []schema.Attr) bool {
	return insql.Has(attrs, &mysql.AutoIncrement{})
}

func findIndexType(attrs []schema.Attr) string {
	var t mysql.IndexType
	if insql.Has(attrs, &t) && t.T != "" {
		return t.T
	} else {
		return "BTREE"
	}
}
