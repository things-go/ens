package driver

import "github.com/things-go/ens"

type Driver interface {
	GetSchema() (ens.Schemaer, error)
	GetEntityMetadata() ([]ens.EntityMetadata, error)
	GetEntity(tb ens.EntityMetadata) (ens.MixinEntity, error)
	GetEntityDefinition(tbName string) (string, error)
}
