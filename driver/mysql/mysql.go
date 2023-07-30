package mysql

import (
	"context"

	_ "ariga.io/atlas/sql/mysql"
	_ "github.com/go-sql-driver/mysql"

	"ariga.io/atlas/sql/schema"
	"ariga.io/atlas/sql/sqlclient"
	"github.com/things-go/ens"
	"github.com/things-go/ens/driver"
)

var _ driver.Driver = (*MySQL)(nil)

type MySQL struct {
	// URL See: https://atlasgo.io/url
	URL string
}

func (self *MySQL) InspectSchema(ctx context.Context, opts *schema.InspectOptions) (ens.Schemaer, error) {
	client, err := sqlclient.Open(ctx, self.URL)
	if err != nil {
		return nil, err
	}
	schemaes, err := client.InspectSchema(ctx, "", opts)
	if err != nil {
		return nil, err
	}
	entities := make([]ens.MixinEntity, 0, len(schemaes.Tables))
	for _, tb := range schemaes.Tables {
		entities = append(entities, IntoEntity(tb))
	}
	return &ens.MixinSchema{
		Name:     schemaes.Name,
		Entities: entities,
	}, nil
}
