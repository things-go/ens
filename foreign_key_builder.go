package ens

import (
	"ariga.io/atlas/sql/schema"
	"github.com/things-go/ens/internal/insql"
)

// ForeignKeyFromDef returns a new ForeignKey with the ForeignKeyDef.
func ForeignKey(symbol string) *foreignKeyBuilder {
	return &foreignKeyBuilder{
		inner: &ForeignKeyDescriptor{
			Symbol:     symbol,
			Table:      "",
			Columns:    nil,
			RefTable:   "",
			RefColumns: nil,
			OnUpdate:   schema.Restrict,
			OnDelete:   schema.Restrict,
			ForeignKey: nil,
		},
	}
}

// ForeignKeyFromDef returns a new ForeignKey with the ForeignKeyDef.
func ForeignKeyFromDef(def ForeignKeyDef) *foreignKeyBuilder {
	fk := def.ForeignKey()
	return &foreignKeyBuilder{
		inner: &ForeignKeyDescriptor{
			Symbol:     fk.Symbol,
			Table:      fk.Table.Name,
			Columns:    insql.ColumnNames(fk.Columns),
			RefTable:   fk.RefTable.Name,
			RefColumns: insql.ColumnNames(fk.RefColumns),
			OnUpdate:   fk.OnUpdate,
			OnDelete:   fk.OnDelete,
			ForeignKey: def,
		},
	}
}

// foreignKeyBuilder is the builder for ForeignKey.
type foreignKeyBuilder struct {
	inner *ForeignKeyDescriptor
}

func (b *foreignKeyBuilder) Table(tbName string, columns []string) *foreignKeyBuilder {
	b.inner.Table = tbName
	b.inner.Columns = columns
	return b
}
func (b *foreignKeyBuilder) RefTable(tbName string, columns []string) *foreignKeyBuilder {
	b.inner.RefTable = tbName
	b.inner.RefColumns = columns
	return b
}
func (b *foreignKeyBuilder) OnDelete(v schema.ReferenceOption) *foreignKeyBuilder {
	b.inner.OnDelete = v
	return b
}
func (b *foreignKeyBuilder) OnUpdate(v schema.ReferenceOption) *foreignKeyBuilder {
	b.inner.OnUpdate = v
	return b
}

// Build implements the ForeignKeyer interface by returning its descriptor.
func (b *foreignKeyBuilder) Build() *ForeignKeyDescriptor {
	return b.inner
}
