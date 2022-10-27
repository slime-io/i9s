package render

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"strings"
)

type ProxyInfoEx struct{}

func (i ProxyInfoEx) IsGeneric() bool {
	return false
}

func (i ProxyInfoEx) Render(o interface{}, ns string, r *Row) error {

	res, ok := o.(ProxyInfoExRes)
	if !ok {
		return fmt.Errorf("expected ProxyInfoEx, but got %T", o)
	}
	id := strings.Join([]string{res.Parent, res.Name}, "#")
	r.ID, r.Fields = id, append(r.Fields, res.Name)
	return nil
}

func (i ProxyInfoEx) Header(ns string) Header {
	return Header{
		HeaderColumn{Name: "Ex"},
	}
}

func (i ProxyInfoEx) ColorerFunc() ColorerFunc {
	return func(ns string, _ Header, re RowEvent) tcell.Color {
		return tcell.ColorCadetBlue
	}
}

type ProxyInfoExRes struct {
	Name   string
	Parent string
}

// GetObjectKind returns a schema object.
func (ProxyInfoExRes) GetObjectKind() schema.ObjectKind {
	return nil
}

// DeepCopyObject returns a container copy.
func (h ProxyInfoExRes) DeepCopyObject() runtime.Object {
	return h
}
