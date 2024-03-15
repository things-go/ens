package codegen

import (
	"fmt"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/things-go/ens"
)

type CodeGenRapier struct {
	OutputDir string
	opts      []Option
}

func NewCodeGenRapier(outputDir string, opts ...Option) *CodeGenRapier {
	return &CodeGenRapier{
		OutputDir: outputDir,
		opts:      opts,
	}
}

func (g *CodeGenRapier) Gen(modelImportPath string, models ...any) error {
	for _, model := range models {
		mixinEntity, err := ens.ParseModel(model)
		if err != nil {
			return err
		}
		entity := mixinEntity.Build(nil)
		data, err := New([]*ens.EntityDescriptor{entity}, g.opts...).
			GenRapier(modelImportPath).
			FormatSource()
		if err != nil {
			return fmt.Errorf("%v: %v", entity.Name, err)
		}
		filename := joinFilename(g.OutputDir, entity.Name, ".rapier.gen.go")
		err = WriteFile(filename, data)
		if err != nil {
			return fmt.Errorf("%v: %v", entity.Name, err)
		}
		slog.Info("👉 " + filename)
	}
	return nil
}

func joinFilename(dir, filename, suffix string) string {
	suffix = strings.TrimSpace(suffix)
	if suffix != "" && !strings.HasPrefix(suffix, ".") {
		suffix = "." + suffix
	}
	return filepath.Join(dir, filename) + suffix
}

// WriteFile writes data to a file named by filename.
// If the file does not exist, WriteFile creates it
// and its upper level paths.
func WriteFile(filename string, data []byte) error {
	if err := os.MkdirAll(path.Dir(filename), os.ModePerm); err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0655)
}
