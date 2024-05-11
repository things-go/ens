package mysql

import (
	"context"

	"ariga.io/atlas/sql/schema"
	"ariga.io/atlas/sql/sqlclient"

	"github.com/things-go/ens"
	"github.com/things-go/ens/driver"
	"github.com/things-go/ens/sqlx"

	_ "ariga.io/atlas/sql/mysql"
	_ "github.com/go-sql-driver/mysql"
)

var _ driver.Driver = (*MySQL)(nil)

type MySQL struct{}

func (self *MySQL) InspectSchema(ctx context.Context, arg *driver.InspectOption) (*ens.Schema, error) {
	schemaes, err := self.inspectSchema(ctx, arg)
	if err != nil {
		return nil, err
	}
	entities := make([]*ens.EntityDescriptor, 0, len(schemaes.Tables))
	for _, tb := range schemaes.Tables {
		entities = append(entities, intoSchema(tb))
	}
	return &ens.Schema{
		Name:     schemaes.Name,
		Entities: entities,
	}, nil
}

// InspectSql implements driver.Driver.
func (self *MySQL) InspectSql(ctx context.Context, arg *driver.InspectOption) (*sqlx.Schema, error) {
	schemaes, err := self.inspectSchema(ctx, arg)
	if err != nil {
		return nil, err
	}
	entities := make([]*sqlx.Table, 0, len(schemaes.Tables))
	for _, tb := range schemaes.Tables {
		entities = append(entities, intoSql(tb))
	}
	return &sqlx.Schema{
		Name:     schemaes.Name,
		Entities: entities,
	}, nil
}

func (self *MySQL) inspectSchema(ctx context.Context, arg *driver.InspectOption) (*schema.Schema, error) {
	client, err := sqlclient.Open(ctx, arg.URL)
	if err != nil {
		return nil, err
	}
	return client.InspectSchema(ctx, "", &arg.InspectOptions)
}
