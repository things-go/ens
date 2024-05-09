package command

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"ariga.io/atlas/sql/schema"
	"github.com/spf13/cobra"
	"github.com/things-go/ens/driver"
	"github.com/things-go/ens/rapier"
)

type rapierOpt struct {
	// sql file
	InputFile []string
	Schema    string
	// database url
	Url     string
	Tables  []string
	Exclude []string

	// output directory
	OutputDir string

	// codegen
	PackageName       string // required, proto 包名
	ModelImportPath   string // required, model导入路径
	DisableDocComment bool   // 禁用doc注释
	EnableInt         bool   // 使能int8,uint8,int16,uint16,int32,uint32输出为int,uint
	EnableIntegerInt  bool   // 使能int32,uint32输出为int,uint
	EnableBoolInt     bool   // 使能bool输出int
}

type rapierCmd struct {
	cmd *cobra.Command
	rapierOpt
}

func newRapierCmd() *rapierCmd {
	root := &rapierCmd{}

	rapierSchema := func() (*rapier.Schema, error) {
		if root.Url != "" {
			d, err := LoadDriver(root.Url)
			if err != nil {
				return nil, err
			}
			return d.InspectRapier(context.Background(), &driver.InspectOption{
				URL: root.Url,
				InspectOptions: schema.InspectOptions{
					Mode:    schema.InspectTables,
					Tables:  root.Tables,
					Exclude: root.Exclude,
				},
			})
		}
		if len(root.InputFile) > 0 {
			d, err := driver.LoadDriver(root.Schema)
			if err != nil {
				return nil, err
			}
			schemas := &rapier.Schema{
				Name:     "",
				Entities: make([]*rapier.Struct, 0, 128),
			}
			for _, filename := range root.InputFile {
				tmpSchema, err := func() (*rapier.Schema, error) {
					content, err := os.ReadFile(filename)
					if err != nil {
						return nil, err
					}
					return d.InspectRapier(context.Background(), &driver.InspectOption{
						URL:            "",
						Data:           string(content),
						InspectOptions: schema.InspectOptions{},
					})
				}()
				if err != nil {
					slog.Warn("🧐 parse failed !!!", slog.String("file", filename), slog.Any("error", err))
					continue
				}
				schemas.Entities = append(schemas.Entities, tmpSchema.Entities...)
			}
			return schemas, nil
		}
		return nil, errors.New("at least one of [url input] is required")
	}

	cmd := &cobra.Command{
		Use:     "rapier",
		Short:   "Generate rapier from database/file",
		Example: "ormat rapier",
		RunE: func(*cobra.Command, []string) error {
			sc, err := rapierSchema()
			if err != nil {
				return err
			}
			for _, msg := range sc.Entities {
				codegen := &rapier.CodeGen{
					Entities:          []*rapier.Struct{msg},
					ByName:            "ormat",
					Version:           version,
					PackageName:       root.PackageName,
					ModelImportPath:   root.ModelImportPath,
					DisableDocComment: root.DisableDocComment,
					EnableInt:         root.EnableInt,
					EnableIntegerInt:  root.EnableIntegerInt,
					EnableBoolInt:     root.EnableBoolInt,
				}

				data, err := codegen.Gen().FormatSource()
				if err != nil {
					return err
				}
				filename := joinFilename(root.OutputDir, msg.TableName, ".rapier.gen.go")
				err = WriteFile(filename, data)
				if err != nil {
					return fmt.Errorf("%v: %w", msg.TableName, err)
				}
				slog.Info("👉 " + filename)
			}
			return nil
		},
	}

	cmd.Flags().StringSliceVarP(&root.InputFile, "input", "i", nil, "input file")
	cmd.Flags().StringVarP(&root.Schema, "schema", "s", "file+mysql", "parser file driver, [file+mysql,file+tidb](仅input时有效)")

	// database url
	cmd.Flags().StringVarP(&root.Url, "url", "u", "", "mysql://root:123456@127.0.0.1:3306/test")
	cmd.Flags().StringSliceVarP(&root.Tables, "table", "t", nil, "only out custom table(仅url时有效)")
	cmd.Flags().StringSliceVarP(&root.Exclude, "exclude", "e", nil, "exclude table pattern(仅url时有效)")

	cmd.Flags().StringVarP(&root.OutputDir, "out", "o", "./repository", "out directory")

	cmd.Flags().StringVar(&root.PackageName, "package", "repository", "proto package name")
	cmd.Flags().StringVar(&root.ModelImportPath, "modelImportPath", "", "model导入路径")
	cmd.Flags().BoolVar(&root.DisableDocComment, "enableInt", false, "禁用文档注释")
	cmd.Flags().BoolVar(&root.EnableInt, "disableBool", false, "使能int8,uint8,int16,uint16,int32,uint32输出为int,uint")
	cmd.Flags().BoolVar(&root.EnableIntegerInt, "enableIntegerInt", false, "使能int32,uint32输出为int,uint")
	cmd.Flags().BoolVar(&root.EnableBoolInt, "enableBoolInt", false, "使能bool输出int")

	cmd.MarkFlagsOneRequired("url", "input")

	root.cmd = cmd
	return root
}
