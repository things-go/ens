package driver

import (
	"context"

	"ariga.io/atlas/sql/schema"
	"github.com/things-go/ens"
)

type Driver interface {
	InspectSchema(context.Context, *schema.InspectOptions) (ens.Schemaer, error)
}
