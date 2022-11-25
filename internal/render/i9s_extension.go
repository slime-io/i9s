package render

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"strings"
)

type I9sExtension struct{}

func (i I9sExtension) IsGeneric() bool {
	return false
}

func (i I9sExtension) Render(o interface{}, ns string, r *Row) error {

	res, ok := o.(I9sExtensionRes)
	if !ok {
		return fmt.Errorf("expected I9sExtensionRes, but got %T", o)
	}
	id := strings.Join([]string{res.Parent, res.Name}, "#")
	r.ID, r.Fields = id, append(r.Fields, res.Name)
	return nil
}

func (i I9sExtension) Header(ns string) Header {
	return Header{
		HeaderColumn{Name: "EXTENSION COMMAND"},
	}
}

func (i I9sExtension) ColorerFunc() ColorerFunc {
	return func(ns string, _ Header, re RowEvent) tcell.Color {
		return tcell.ColorCadetBlue
	}
}

type I9sExtensionRes struct {
	Name   string
	Rev    string
	Parent string
}

// GetObjectKind returns a schema object.
func (i I9sExtensionRes) GetObjectKind() schema.ObjectKind {
	return nil
}

// DeepCopyObject returns a container copy.
func (i I9sExtensionRes) DeepCopyObject() runtime.Object {
	return i
}
