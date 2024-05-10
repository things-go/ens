package ens

import "github.com/things-go/ens/internal/insql"

var _ MixinEntity = (*EntityBuilder)(nil)

type EntityBuilder struct {
	name        string         // schema entity name
	comment     string         // schema entity comment
	table       TableDef       // entity table define
	fields      []Fielder      // field information
	indexes     []Indexer      // index information
	foreignKeys []ForeignKeyer // foreign key information
}

// EntityFromDef returns a new entity with the TableDef.
// auto set name, comment, table.
func EntityFromDef(def TableDef) *EntityBuilder {
	tb := def.Table()
	return &EntityBuilder{
		name:    tb.Name,
		comment: insql.MustComment(tb.Attrs),
		table:   def,
		fields:  nil,
		indexes: nil,
	}
}

func (self *EntityBuilder) SetMetadata(name, comment string) *EntityBuilder {
	self.name = name
	self.comment = comment
	return self
}
func (self *EntityBuilder) SetTable(tb TableDef) *EntityBuilder {
	self.table = tb
	return self
}
func (self *EntityBuilder) SetFields(fields ...Fielder) *EntityBuilder {
	self.fields = fields
	return self
}
func (self *EntityBuilder) SetIndexes(indexes ...Indexer) *EntityBuilder {
	self.indexes = indexes
	return self
}
func (self *EntityBuilder) SetForeignKeys(fks ...ForeignKeyer) *EntityBuilder {
	self.foreignKeys = fks
	return self
}
func (self *EntityBuilder) Metadata() (name, comment string) { return self.name, self.comment }
func (self *EntityBuilder) Table() TableDef                  { return self.table }
func (self *EntityBuilder) Fields() []Fielder                { return self.fields }
func (self *EntityBuilder) Indexes() []Indexer               { return self.indexes }
func (self *EntityBuilder) ForeignKeys() []ForeignKeyer      { return self.foreignKeys }
func (self *EntityBuilder) Build(opt *Option) *EntityDescriptor {
	fielders := self.Fields()
	fields := make([]*FieldDescriptor, 0, len(fielders))
	for _, fb := range fielders {
		field := fb.Build(opt)
		fields = append(fields, field)
	}
	indexers := self.Indexes()
	indexes := make([]*IndexDescriptor, 0, len(indexers))
	for _, v := range indexers {
		indexes = append(indexes, v.Build())
	}
	fkers := self.ForeignKeys()
	fks := make([]*ForeignKeyDescriptor, 0, len(fkers))
	for _, v := range fkers {
		fks = append(fks, v.Build())
	}
	name, comment := self.Metadata()
	return &EntityDescriptor{
		Name:        name,
		Comment:     comment,
		Table:       self.Table(),
		Fields:      fields,
		Indexes:     indexes,
		ForeignKeys: fks,
	}
}
