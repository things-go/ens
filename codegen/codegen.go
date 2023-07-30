package codegen

import (
	"bytes"
	"fmt"

	"github.com/things-go/ens"
	"golang.org/x/tools/imports"
)

type CodeGen struct {
	buf               bytes.Buffer
	entities          []*ens.EntityDescriptor
	byName            string
	version           string
	packageName       string
	options           map[string]string
	skipColumns       map[string]struct{}
	hasColumn         bool
	disableDocComment bool
}

type Option func(*CodeGen)

func WithByName(s string) Option {
	return func(g *CodeGen) {
		g.byName = s
	}
}

func WithVersion(version string) Option {
	return func(g *CodeGen) {
		g.version = version
	}
}

func WithPackageName(s string) Option {
	return func(g *CodeGen) {
		g.packageName = s
	}
}

func WithOptions(mp map[string]string) Option {
	return func(g *CodeGen) {
		for k, v := range mp {
			g.options[k] = v
		}
	}
}

func WithSkipColumns(mp map[string]struct{}) Option {
	return func(g *CodeGen) {
		if mp != nil {
			g.skipColumns = mp
		}
	}
}

func WithHasColumn(b bool) Option {
	return func(g *CodeGen) {
		g.hasColumn = b
	}
}

func WithDisableDocComment(b bool) Option {
	return func(g *CodeGen) {
		g.disableDocComment = b
	}
}

func New(md []*ens.EntityDescriptor, opts ...Option) *CodeGen {
	g := &CodeGen{
		entities:    md,
		byName:      "codegen",
		version:     "unknown",
		packageName: "codegen",
		options:     make(map[string]string),
		skipColumns: make(map[string]struct{}),
		hasColumn:   false,
	}
	for _, opt := range opts {
		opt(g)
	}
	return g
}

// Bytes returns the CodeBuf's buffer.
func (g *CodeGen) Bytes() []byte {
	return g.buf.Bytes()
}

// FormatSource return formats and adjusts imports contents of the CodeGen's buffer.
func (g *CodeGen) FormatSource() ([]byte, error) {
	data := g.buf.Bytes()
	if len(data) == 0 {
		return data, nil
	}
	// return format.Source(data)
	return imports.Process("", data, nil)
}

// Write appends the contents of p to the buffer,
func (g *CodeGen) Write(b []byte) (n int, err error) {
	return g.buf.Write(b)
}

// P prints a line to the generated output. It converts each parameter to a
// string following the same rules as fmt.Print. It never inserts spaces
// between parameters.
func (g *CodeGen) P(args ...any) *CodeGen {
	for _, arg := range args {
		fmt.Fprint(&g.buf, arg)
	}
	fmt.Fprintln(&g.buf)
	return g
}
