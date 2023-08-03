package driver

import (
	"context"
	"sync"

	"ariga.io/atlas/sql/schema"
	"github.com/things-go/ens"
)

const (
	Mysql         = "mysql"
	FileMysql     = "file+mysql"
	FileMysqlTidb = "file+tidb"
)

var drivers sync.Map

type Driver interface {
	InspectSchema(context.Context, *InspectOption) (ens.Schemaer, error)
}

func RegisterDriver(name string, d Driver) {
	if _, ok := drivers.Load(name); ok {
		panic("driver: Register called twice for " + name)
	}
	drivers.Store(name, d)
}

func LoadDriver(name string) (Driver, bool) {
	d, ok := drivers.Load(name)
	if !ok {
		return nil, false
	}
	return d.(Driver), true
}

func DriverNames() []string {
	var names []string

	drivers.Range(func(key, value any) bool {
		names = append(names, key.(string))
		return true
	})
	return names
}

type InspectOption struct {
	// URL See: https://atlasgo.io/url
	URL string
	// sql data, for file
	Data string
	// InspectOptions describes options for Inspector.
	schema.InspectOptions
}
