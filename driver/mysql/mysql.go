package mysql

import (
	"context"

	_ "ariga.io/atlas/sql/mysql"
	_ "github.com/go-sql-driver/mysql"

	"ariga.io/atlas/sql/sqlclient"
	"github.com/things-go/ens"
	"github.com/things-go/ens/driver"
)

var _ driver.Driver = (*MySQL)(nil)

type MySQL struct{}

func (self *MySQL) InspectSchema(ctx context.Context, arg *driver.InspectOption) (*ens.MixinSchema, error) {
	client, err := sqlclient.Open(ctx, arg.URL)
	if err != nil {
		return nil, err
	}
	schemaes, err := client.InspectSchema(ctx, "", &arg.InspectOptions)
	if err != nil {
		return nil, err
	}
	entities := make([]ens.MixinEntity, 0, len(schemaes.Tables))
	for _, tb := range schemaes.Tables {
		entities = append(entities, IntoMixinEntity(tb))
	}
	return &ens.MixinSchema{
		Name:     schemaes.Name,
		Entities: entities,
	}, nil
}
