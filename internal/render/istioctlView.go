package render

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"strings"
)

type IstioctlView struct{}

func (i IstioctlView) IsGeneric() bool {
	return false
}

func (i IstioctlView) Render(o interface{}, ns string, r *Row) error {

	res, ok := o.(IstioctlViewRes)
	if !ok {
		return fmt.Errorf("expected EnvoyAPIRes, but got %T", o)
	}
	id := strings.Join([]string{res.Parent, res.Name}, "#")
	r.ID, r.Fields = id, append(r.Fields, res.Name)
	return nil
}

func (i IstioctlView) Header(ns string) Header {
	return Header{
		HeaderColumn{Name: "istioctl command items"},
	}
}

func (i IstioctlView) ColorerFunc() ColorerFunc {
	return func(ns string, _ Header, re RowEvent) tcell.Color {
		return tcell.ColorCadetBlue
	}
}

type IstioctlViewRes struct {
	Name   string
	Parent string
}

// GetObjectKind returns a schema object.
func (IstioctlViewRes) GetObjectKind() schema.ObjectKind {
	return nil
}

// DeepCopyObject returns a container copy.
func (h IstioctlViewRes) DeepCopyObject() runtime.Object {
	return h
}
