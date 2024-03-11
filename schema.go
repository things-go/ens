package ens

// Schema
type Schema struct {
	Name     string              // schema name
	Entities []*EntityDescriptor // schema entity.
}

// MixinSchema information
type MixinSchema struct {
	Name     string        // schema name
	Entities []MixinEntity // schema entity.
}

func (self *MixinSchema) Build(opt *Option) *Schema {
	entities := make([]*EntityDescriptor, 0, len(self.Entities))
	for _, mixin := range self.Entities {
		entities = append(entities, mixin.Build(opt))
	}
	return &Schema{
		Name:     self.Name,
		Entities: entities,
	}
}
