package ens

import (
	"github.com/things-go/ens/internal/sqlx"
)

// EntityDescriptor Each table corresponds to an EntityDescriptor
type EntityDescriptor struct {
	Name         string                  // entity name
	Comment      string                  // entity comment
	Table        TableDef                // entity table define
	Fields       []*FieldDescriptor      // field information
	Indexes      []*IndexDescriptor      // index information
	ForeignKeys  []*ForeignKeyDescriptor // foreign key information
	ProtoMessage []*ProtoMessage         // protobuf message information.
}

type EntityDescriptorSlice []*EntityDescriptor

func (t EntityDescriptorSlice) Len() int           { return len(t) }
func (t EntityDescriptorSlice) Less(i, j int) bool { return t[i].Name < t[j].Name }
func (t EntityDescriptorSlice) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }

func BuildEntity(m MixinEntity, opt *Option) *EntityDescriptor {
	enableGogo, enableSea := false, false
	if opt != nil {
		enableGogo, enableSea = opt.EnableGogo, opt.EnableSea
	}
	fielders := m.Fields()
	fields := make([]*FieldDescriptor, 0, len(fielders))
	protoMessages := make([]*ProtoMessage, 0, len(fielders))
	for _, fb := range fielders {
		field := fb.Build(opt)
		fields = append(fields, field)
		protoMessages = append(protoMessages, buildProtoMessage(field, enableGogo, enableSea))
	}
	indexers := m.Indexes()
	indexes := make([]*IndexDescriptor, 0, len(indexers))
	for _, v := range indexers {
		indexes = append(indexes, v.Build())
	}
	fkers := m.ForeignKeys()
	fks := make([]*ForeignKeyDescriptor, 0, len(fkers))
	for _, v := range fkers {
		fks = append(fks, v.Build())
	}
	name, comment := m.Metadata()
	return &EntityDescriptor{
		Name:         name,
		Comment:      comment,
		Table:        m.Table(),
		Fields:       fields,
		Indexes:      indexes,
		ForeignKeys:  fks,
		ProtoMessage: protoMessages,
	}
}

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
		comment: sqlx.MustComment(tb.Attrs),
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
