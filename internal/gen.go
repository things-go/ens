// gen is a codegen cmd for generating numeric build types from template.
package main

import (
	"bytes"
	"embed"
	"go/format"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/things-go/ens"
)

//go:embed integer.tpl float.tpl
var staticFile embed.FS

var components = template.Must(
	template.New("components").
		Funcs(template.FuncMap{
			"title":     strings.Title,
			"hasPrefix": strings.HasPrefix,
			"toUpper":   strings.ToUpper,
		}).
		ParseFS(staticFile, "*.tpl"),
)

type Metadata struct {
	Kind ens.Type
}

type Numeric struct {
	tpl        *template.Template
	fieldTypes []ens.Type
}

func main() {
	numerics := []Numeric{
		{
			tpl: components.Lookup("integer.tpl"),
			fieldTypes: []ens.Type{
				ens.TypeInt,
				ens.TypeUint,
				ens.TypeInt8,
				ens.TypeInt16,
				ens.TypeInt32,
				ens.TypeInt64,
				ens.TypeUint8,
				ens.TypeUint16,
				ens.TypeUint32,
				ens.TypeUint64,
			},
		},
		{
			tpl: components.Lookup("float.tpl"),
			fieldTypes: []ens.Type{
				ens.TypeFloat64,
				ens.TypeFloat32,
			},
		},
	}

	b := &bytes.Buffer{}
	for _, v := range numerics {
		for _, t := range v.fieldTypes {
			b.Reset()
			err := v.tpl.Execute(b, Metadata{Kind: t})
			if err != nil {
				log.Fatal("executing template:", err)
			}
			buf, err := format.Source(b.Bytes())
			if err != nil {
				log.Fatal("formatting output:", err)
			}
			filename := "field_" + t.String() + ".go"
			if err := os.WriteFile(filename, buf, 0644); err != nil {
				log.Fatal("writing go file:", err)
			}
		}
	}
}
