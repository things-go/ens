package command

import (
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/things-go/ens/sqlx"
)

type sqlOpt struct {
	URL               string
	Tables            []string
	Exclude           []string
	DisableDocComment bool

	OutputDir string
	Merge     bool
	Filename  string
}

type sqlCmd struct {
	cmd *cobra.Command
	sqlOpt
}

// must be same as atlas schema inspect -u "mysql://localhost" --format "{{ sql . }}"
func newSqlCmd() *sqlCmd {
	root := &sqlCmd{}
	cmd := &cobra.Command{
		Use:     "sql",
		Short:   "Generate sql file",
		Example: "ormat sql",
		RunE: func(*cobra.Command, []string) error {
			sc, err := getSchema(&source{
				URL:     root.URL,
				Tables:  root.Tables,
				Exclude: root.Exclude,
			})
			if err != nil {
				return err
			}
			schemaes := sc.IntoSQL()
			if root.Merge {
				codegen := &sqlx.CodeGen{
					Entities:          schemaes.Entities,
					ByName:            "ormat",
					Version:           version,
					DisableDocComment: root.DisableDocComment,
				}
				data := codegen.Gen().Bytes()
				filename := joinFilename(root.OutputDir, root.Filename, ".sql")
				err = WriteFile(filename, data)
				if err != nil {
					return err
				}
				slog.Info("👉 " + filename)
			} else {
				for _, entity := range schemaes.Entities {
					codegen := &sqlx.CodeGen{
						Entities:          []*sqlx.Table{entity},
						ByName:            "ormat",
						Version:           version,
						DisableDocComment: root.DisableDocComment,
					}
					data := codegen.Gen().Bytes()
					filename := joinFilename(root.OutputDir, entity.Name, ".sql")
					err = WriteFile(filename, data)
					if err != nil {
						return err
					}
					slog.Info("👉 " + filename)
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&root.URL, "url", "", "mysql://root:123456@127.0.0.1:3306/test)")
	cmd.PersistentFlags().StringSliceVarP(&root.Tables, "table", "t", nil, "only out custom table")
	cmd.PersistentFlags().StringSliceVar(&root.Exclude, "exclude", nil, "exclude table pattern")
	cmd.Flags().StringVarP(&root.OutputDir, "out", "o", "./migration", "out directory")
	cmd.Flags().StringVar(&root.Filename, "filename", "create_table", "filename when merge enabled")
	cmd.Flags().BoolVar(&root.Merge, "merge", false, "merge in a file")
	cmd.Flags().BoolVarP(&root.DisableDocComment, "disableDocComment", "d", false, "禁用文档注释")

	cmd.MarkFlagRequired("url")

	root.cmd = cmd
	return root
}
