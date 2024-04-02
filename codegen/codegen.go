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

// WithByName the code generator by which executables name.
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

// Print formats using the default formats for its operands and writes to the generated output.
// Spaces are added between operands when neither is a string.
// It returns the number of bytes written and any write error encountered.
func (g *CodeGen) Print(a ...any) (n int, err error) {
	return fmt.Fprint(&g.buf, a...)
}

// Printf formats according to a format specifier for its operands and writes to the generated output.
// It returns the number of bytes written and any write error encountered.
func (g *CodeGen) Printf(format string, a ...any) (n int, err error) {
	return fmt.Fprintf(&g.buf, format, a...)
}

// Fprintln formats using the default formats to the generated output.
// Spaces are always added between operands and a newline is appended.
// It returns the number of bytes written and any write error encountered.
func (g *CodeGen) Println(a ...any) (n int, err error) {
	return fmt.Fprintln(&g.buf, a...)
}
