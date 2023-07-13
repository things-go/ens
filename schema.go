package ens

// Schema
type Schema struct {
	Name     string    // schema name
	Entities []*Entity // schema entity.
}

// MixinSchema information
type MixinSchema struct {
	Name     string        // schema name
	Entities []MixinEntity // schema entity.
}

func (self *MixinSchema) Build(opt *Option) *Schema {
	entities := make([]*Entity, 0, len(self.Entities))
	for _, mixin := range self.Entities {
		entities = append(entities, BuildEntity(mixin, opt))
	}
	return &Schema{
		Name:     self.Name,
		Entities: entities,
	}
}
