package ens

import (
	"github.com/things-go/ens/proto"
	"github.com/things-go/ens/rapier"
)

// Schema
type Schema struct {
	Name     string              // schema name
	Entities []*EntityDescriptor // schema entity.
}

func (s *Schema) IntoProto() *proto.Schema {
	entities := make([]*proto.Message, 0, len(s.Entities))
	for _, entity := range s.Entities {
		entities = append(entities, entity.IntoProto())
	}
	return &proto.Schema{
		Name:     s.Name,
		Entities: entities,
	}
}

func (s *Schema) IntoRapier() *rapier.Schema {
	entities := make([]*rapier.Struct, 0, len(s.Entities))
	for _, entity := range s.Entities {
		entities = append(entities, entity.IntoRapier())
	}
	return &rapier.Schema{
		Name:     s.Name,
		Entities: entities,
	}
}
