package insql

import (
	"fmt"
	"reflect"

	"ariga.io/atlas/sql/schema"
)

var (
	attrsType   = reflect.TypeOf(([]schema.Attr)(nil))
	clausesType = reflect.TypeOf(([]schema.Clause)(nil))
	exprsType   = reflect.TypeOf(([]schema.Expr)(nil))
)

// Has finds the first element in the elements list that
// matches target, and if so, sets target to that attribute
// value and returns true.
// support schema.Attr, schema.Clause, schema.Expr
func Has(elements, target any) bool {
	ev := reflect.ValueOf(elements)
	if t := ev.Type(); t != attrsType && t != clausesType && t != exprsType {
		panic(fmt.Sprintf("unexpected elements type: %T", elements))
	}
	tv := reflect.ValueOf(target)
	if tv.Kind() != reflect.Ptr || tv.IsNil() {
		panic("target must be a non-nil pointer")
	}
	for i := 0; i < ev.Len(); i++ {
		idx := ev.Index(i)
		if idx.IsNil() {
			continue
		}
		if e := idx.Elem(); e.Type().AssignableTo(tv.Type()) {
			tv.Elem().Set(e.Elem())
			return true
		}
	}
	return false
}

func Comment(elements []schema.Attr) (string, bool) {
	var comment schema.Comment
	ok := Has(elements, &comment)
	return comment.Text, ok
}

func MustComment(elements []schema.Attr) string {
	v, _ := Comment(elements)
	return v
}

func Charset(elements []schema.Attr) (string, bool) {
	var val schema.Charset
	ok := Has(elements, &val)
	return val.V, ok
}

func Collation(elements []schema.Attr) (string, bool) {
	var val schema.Collation
	ok := Has(elements, &val)
	return val.V, ok
}

func IndexEqual(idx1, idx2 *schema.Index) bool {
	return idx1 != nil && idx2 != nil && (idx1 == idx2 || idx1.Name == idx2.Name)
}

func FindIndexPartSeq(parts []*schema.IndexPart, col *schema.Column) (int, bool) {
	for _, p := range parts {
		if p.C == col || p.C.Name == col.Name {
			return p.SeqNo, true
		}
	}
	return 0, false
}

func IndexPartColumnNames(parts []*schema.IndexPart) []string {
	fields := make([]string, 0, len(parts))
	for _, v := range parts {
		fields = append(fields, v.C.Name)
	}
	return fields
}

func FindColumn(columns []*schema.Column, columnName string) (*schema.Column, bool) {
	for _, col := range columns {
		if col.Name == columnName {
			return col, true
		}
	}
	return nil, false
}

func ColumnNames(columns []*schema.Column) []string {
	ns := make([]string, 0, len(columns))
	for _, col := range columns {
		ns = append(ns, col.Name)
	}
	return ns
}

// P returns a pointer to v.
func P[T any](v T) *T {
	return &v
}

// V returns the value p is pointing to.
// If p is nil, the zero value is returned.
func V[T any](p *T) (v T) {
	if p != nil {
		v = *p
	}
	return
}
