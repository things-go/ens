package command

import (
	"cmp"
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/things-go/ens"
	"github.com/things-go/ens/utils"
)

type modelOpt struct {
	source

	OutputDir string

	PackageName string // 包名

	ens.Option
	DisableCommentTag bool // 禁用注释放入tag标签中
	DisableDocComment bool // 禁用文档注释

	Merge         bool
	MergeFilename string
}

type modelCmd struct {
	cmd *cobra.Command
	modelOpt
}

func newModelCmd() *modelCmd {
	root := &modelCmd{}

	cmd := &cobra.Command{
		Use:     "model",
		Short:   "Generate model from database",
		Example: "ormat model",
		RunE: func(*cobra.Command, []string) error {
			schemaes, err := getSchema(&root.source)
			if err != nil {
				return err
			}
			if root.Merge {
				g := ens.CodeGen{
					Entities:          schemaes.Entities,
					ByName:            "ormat",
					Version:           version,
					PackageName:       cmp.Or(root.PackageName, utils.GetPkgName(root.OutputDir)),
					DisableDocComment: root.DisableDocComment,
					Option:            root.Option,
				}
				data, err := g.Gen().FormatSource()
				if err != nil {
					return err
				}
				filename := joinFilename(root.OutputDir, root.MergeFilename, ".go")
				err = WriteFile(filename, data)
				if err != nil {
					return err
				}
				slog.Info("👉 " + filename)
			} else {
				for _, entity := range schemaes.Entities {
					g := &ens.CodeGen{
						Entities:          []*ens.EntityDescriptor{entity},
						ByName:            "ormat",
						Version:           version,
						PackageName:       utils.GetPkgName(root.OutputDir),
						DisableDocComment: root.DisableDocComment,
						Option:            root.Option,
					}
					data, err := g.Gen().FormatSource()
					if err != nil {
						return fmt.Errorf("%v: %v", entity.Name, err)
					}
					filename := joinFilename(root.OutputDir, entity.Name, ".go")
					err = WriteFile(filename, data)
					if err != nil {
						return fmt.Errorf("%v: %v", entity.Name, err)
					}
					slog.Info("👉 " + filename)
				}
			}
			slog.Info("😄 generate success !!!")
			return nil
		},
	}
	// input file
	cmd.Flags().StringSliceVarP(&root.InputFile, "input", "i", nil, "input file")
	cmd.Flags().StringVarP(&root.Schema, "schema", "s", "file+mysql", "parser file driver, [file+mysql,file+tidb](仅input时有效)")
	// database url
	cmd.PersistentFlags().StringVar(&root.URL, "url", "", "mysql://root:123456@127.0.0.1:3306/test")
	cmd.PersistentFlags().StringSliceVarP(&root.Tables, "table", "t", nil, "only out custom table")
	cmd.PersistentFlags().StringSliceVar(&root.Exclude, "exclude", nil, "exclude table pattern")

	cmd.PersistentFlags().StringVarP(&root.OutputDir, "out", "o", "./model", "out directory")

	cmd.PersistentFlags().StringVar(&root.PackageName, "package", "", "package name")
	cmd.PersistentFlags().BoolVarP(&root.DisableDocComment, "disableDocComment", "d", false, "禁用文档注释")

	cmd.PersistentFlags().StringToStringVarP(&root.Tags, "tags", "K", map[string]string{"json": utils.StyleSmallCamelCase}, "tags标签,类型支持[smallCamelCase,camelCase,snakeCase,kebab]")
	cmd.PersistentFlags().BoolVarP(&root.EnableInt, "enableInt", "e", false, "使能int8,uint8,int16,uint16,int32,uint32输出为int,uint")
	cmd.PersistentFlags().BoolVarP(&root.EnableBoolInt, "enableBoolInt", "b", false, "使能bool输出int")
	cmd.PersistentFlags().BoolVarP(&root.DisableNullToPoint, "disableNullToPoint", "B", false, "禁用字段为null时输出指针类型,将输出为sql.Nullxx")
	cmd.PersistentFlags().BoolVarP(&root.DisableCommentTag, "disableCommentTag", "j", false, "禁用注释放入tag标签中")
	cmd.PersistentFlags().BoolVarP(&root.EnableForeignKey, "enableForeignKey", "J", false, "使用外键")

	cmd.PersistentFlags().BoolVar(&root.Merge, "merge", false, "merge in a file or not")
	cmd.PersistentFlags().StringVar(&root.MergeFilename, "model", "", "merge filename")

	cmd.MarkPersistentFlagRequired("url") // nolint

	root.cmd = cmd
	return root
}
