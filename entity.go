package ens

// EntityMetadata entity metadata.
type EntityMetadata struct {
	Name       string // entity name
	Comment    string // entity comment
	Definition string // entity SQL definition.
}

// Entity Each table corresponds to an Entity
type Entity struct {
	Name         string             // entity name
	Comment      string             // entity comment
	Definition   string             // entity SQL definition.
	Fields       []*FieldDescriptor // field information
	Indexes      []*IndexDescriptor // index information
	ProtoMessage []*ProtoMessage    // protobuf message information.
}

type EntitySlice []*Entity

func (t EntitySlice) Len() int           { return len(t) }
func (t EntitySlice) Less(i, j int) bool { return t[i].Name < t[j].Name }
func (t EntitySlice) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }

func BuildEntity(m MixinEntity, opt *Option) *Entity {
	fielders := m.Fields()
	fields := make([]*FieldDescriptor, 0, len(fielders))
	protoMessages := make([]*ProtoMessage, 0, len(fielders))
	for _, fb := range fielders {
		field := fb.Build(opt)
		fields = append(fields, field)
		enableGogo, enableSea := false, false
		if opt != nil {
			enableGogo, enableSea = opt.EnableGogo, opt.EnableSea
		}
		protoMessages = append(protoMessages, buildProtoMessage(field, enableGogo, enableSea))
	}
	indexers := m.Indexes()
	indexes := make([]*IndexDescriptor, 0, len(indexers))
	for _, v := range indexers {
		indexes = append(indexes, v.Build())
	}
	md := m.Metadata()
	return &Entity{
		Name:         md.Name,
		Comment:      md.Comment,
		Definition:   md.Definition,
		Fields:       fields,
		Indexes:      indexes,
		ProtoMessage: protoMessages,
	}
}

var _ MixinEntity = (*EntityBuilder)(nil)

type EntityBuilder struct {
	md      EntityMetadata // schema metadata
	fields  []Fielder      // field information
	indexes []Indexer      // index information
}

func (self *EntityBuilder) SetMetadata(md EntityMetadata) *EntityBuilder {
	self.md = md
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

func (self *EntityBuilder) Fields() []Fielder {
	return self.fields
}

func (self *EntityBuilder) Indexes() []Indexer {
	return self.indexes
}

func (self *EntityBuilder) Metadata() EntityMetadata {
	return self.md
}
