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
	"github.com/things-go/ens/proto"
)

type protoOpt struct {
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
	PackageName       string            // required, proto 包名
	Options           map[string]string // required, proto option
	DisableDocComment bool              // 禁用doc注释
	DisableBool       bool              // 禁用bool,使用int32
	DisableTimestamp  bool              // 禁用google.protobuf.Timestamp,使用int64
}

type protoCmd struct {
	cmd *cobra.Command
	protoOpt
}

func newProtoCmd() *protoCmd {
	root := &protoCmd{}

	protoSchema := func() (*proto.Schema, error) {
		if root.Url != "" {
			d, err := LoadDriver(root.Url)
			if err != nil {
				return nil, err
			}
			return d.InspectProto(context.Background(), &driver.InspectOption{
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
			schemas := &proto.Schema{
				Name:     "",
				Messages: make([]*proto.Message, 0, 128),
			}
			for _, filename := range root.InputFile {
				tmpSchema, err := func() (*proto.Schema, error) {
					content, err := os.ReadFile(filename)
					if err != nil {
						return nil, err
					}
					return d.InspectProto(context.Background(), &driver.InspectOption{
						URL:            "",
						Data:           string(content),
						InspectOptions: schema.InspectOptions{},
					})
				}()
				if err != nil {
					slog.Warn("🧐 parse failed !!!", slog.String("file", filename), slog.Any("error", err))
					continue
				}
				schemas.Messages = append(schemas.Messages, tmpSchema.Messages...)
			}
			return schemas, nil
		}
		return nil, errors.New("at least one of [url input] is required")
	}

	cmd := &cobra.Command{
		Use:     "proto",
		Short:   "Generate proto from database",
		Example: "ormat proto",
		RunE: func(*cobra.Command, []string) error {
			sc, err := protoSchema()
			if err != nil {
				return err
			}
			for _, msg := range sc.Messages {
				codegen := &proto.CodeGen{
					Messages:          []*proto.Message{msg},
					ByName:            "ormat",
					Version:           version,
					PackageName:       root.PackageName,
					Options:           root.Options,
					DisableDocComment: root.DisableDocComment,
					DisableBool:       root.DisableBool,
					DisableTimestamp:  root.DisableTimestamp,
				}
				data := codegen.Gen().Bytes()
				filename := joinFilename(root.OutputDir, msg.TableName, ".proto")
				err := WriteFile(filename, data)
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

	cmd.Flags().StringVarP(&root.OutputDir, "out", "o", "./mapper", "out directory")

	cmd.Flags().StringVar(&root.PackageName, "package", "mapper", "proto package name")
	cmd.Flags().StringToStringVar(&root.Options, "options", nil, "proto options key/value")
	cmd.Flags().BoolVar(&root.DisableDocComment, "disableDocComment", false, "禁用文档注释")
	cmd.Flags().BoolVar(&root.DisableBool, "disableBool", false, "禁用bool,使用int32")
	cmd.Flags().BoolVar(&root.DisableTimestamp, "disableTimestamp", false, "禁用google.protobuf.Timestamp,使用int64")

	cmd.MarkFlagsOneRequired("url", "input")

	root.cmd = cmd
	return root
}
