package command

import (
	"fmt"
	"log/slog"

	"github.com/things-go/ens"
	"github.com/things-go/ens/codegen"
	"github.com/things-go/ens/utils"
)

type genFileOpt struct {
	OutputDir     string
	View          Config
	Merge         bool
	MergeFilename string
	Template      string
}

func (self *genFileOpt) build(mixin ens.Schemaer) *ens.Schema {
	return mixin.Build(&self.View.Option)
}

func (self *genFileOpt) GenModel(mixin ens.Schemaer) error {
	sc := self.build(mixin)
	if self.Merge {
		g := codegen.CodeGen{
			Entities:          sc.Entities,
			ByName:            "ormat",
			Version:           version,
			PackageName:       utils.GetPkgName(self.OutputDir),
			DisableDocComment: self.View.DisableDocComment,
		}
		data, err := g.Gen().FormatSource()
		if err != nil {
			return err
		}
		filename := joinFilename(self.OutputDir, self.MergeFilename, ".go")
		err = WriteFile(filename, data)
		if err != nil {
			return err
		}
		slog.Info("👉 " + filename)
	} else {
		for _, entity := range sc.Entities {
			g := codegen.CodeGen{
				Entities:          []*ens.EntityDescriptor{entity},
				ByName:            "ormat",
				Version:           version,
				PackageName:       utils.GetPkgName(self.OutputDir),
				DisableDocComment: self.View.DisableDocComment,
			}
			data, err := g.Gen().FormatSource()
			if err != nil {
				return fmt.Errorf("%v: %v", entity.Name, err)
			}
			filename := joinFilename(self.OutputDir, entity.Name, ".go")
			err = WriteFile(filename, data)
			if err != nil {
				return fmt.Errorf("%v: %v", entity.Name, err)
			}
			slog.Info("👉 " + filename)
		}
	}
	slog.Info("😄 generate success !!!")
	return nil
}
