package mysql

import (
	"ariga.io/atlas/sql/mysql"
	"ariga.io/atlas/sql/schema"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/things-go/ens"
	"github.com/things-go/ens/internal/sqlx"
	"github.com/things-go/ens/proto"
	"github.com/things-go/ens/rapier"
	"github.com/things-go/ens/utils"
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

func IntoProto(tb *schema.Table) *proto.Message {
	// * columns
	fields := make([]*proto.MessageField, 0, len(tb.Columns))
	for _, col := range tb.Columns {
		goType := intoGoType(col.Type.Raw)
		k, n := goType.Type.IntoProtoKind()
		cardinality := protoreflect.Required
		if col.Type.Null {
			cardinality = protoreflect.Optional
		}
		fields = append(fields, &proto.MessageField{
			Cardinality: cardinality,
			Type:        k,
			TypeName:    n,
			Name:        col.Name,
			ColumnName:  col.Name,
			Comment:     sqlx.MustComment(col.Attrs),
		})
	}
	return &proto.Message{
		Name:      tb.Name,
		TableName: tb.Name,
		Comment:   sqlx.MustComment(tb.Attrs),
		Fields:    fields,
	}
}

func IntoRapier(tb *schema.Table) *rapier.Struct {
	// * columns
	fields := make([]*rapier.StructField, 0, len(tb.Columns))
	for _, col := range tb.Columns {
		goType := intoGoType(col.Type.Raw)

		t := goType.Type.IntoRapierType()

		fields = append(fields, &rapier.StructField{
			Type:       t,
			GoName:     utils.CamelCase(col.Name),
			Nullable:   col.Type.Null,
			ColumnName: col.Name,
			Comment:    sqlx.MustComment(col.Attrs),
		})
	}
	return &rapier.Struct{
		GoName:    utils.CamelCase(tb.Name),
		TableName: tb.Name,
		Comment:   sqlx.MustComment(tb.Attrs),
		Fields:    fields,
	}
}
