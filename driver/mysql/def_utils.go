package mysql

import (
	"ariga.io/atlas/sql/mysql"
	"ariga.io/atlas/sql/schema"
	"github.com/things-go/ens"
	"github.com/things-go/ens/internal/sqlx"
)

func autoIncrement(attrs []schema.Attr) bool {
	return sqlx.Has(attrs, &mysql.AutoIncrement{})
}

func findIndexType(attrs []schema.Attr) string {
	var t mysql.IndexType
	if sqlx.Has(attrs, &t) && t.T != "" {
		return t.T
	} else {
		return "BTREE"
	}
}

func IntoMixinEntity(tb *schema.Table) ens.MixinEntity {
	// * columns
	fielders := make([]ens.Fielder, 0, len(tb.Columns))
	for _, col := range tb.Columns {
		colDef := NewColumnDef(col)
		fielders = append(fielders,
			ens.FieldFromDef(intoGoType(col.Type.Raw), colDef).
				Tags(colDef.GormTag(tb)),
		)
	}
	// * indexes
	indexers := make([]ens.Indexer, 0, len(tb.Indexes))
	for _, index := range tb.Indexes {
		indexers = append(indexers, ens.IndexFromDef(NewIndexDef(index)))
	}
	//* foreignKeys
	fkers := make([]ens.ForeignKeyer, 0, len(tb.ForeignKeys))
	for _, fk := range tb.ForeignKeys {
		fkers = append(fkers, ens.ForeignKeyFromDef(NewForeignKey(fk)))
	}

	// * table
	return ens.EntityFromDef(NewTableDef(tb)).
		SetFields(fielders...).
		SetIndexes(indexers...).
		SetForeignKeys(fkers...)
}
