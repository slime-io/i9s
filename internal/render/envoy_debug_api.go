package render

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type EnvoyApi struct{}

func (i EnvoyApi) IsGeneric() bool {
	return false
}

func (i EnvoyApi) Render(o interface{}, ns string, r *Row) error {

	res, ok := o.(EnvoyApiRes)
	if !ok {
		return fmt.Errorf("expected EnvoyAPIRes, but got %T", o)
	}
	r.ID, r.Fields = res.Name, append(r.Fields, res.Name)
	return nil
}

func (i EnvoyApi) Header(ns string) Header {
	return Header{
		HeaderColumn{Name: "api"},
	}
}

func (i EnvoyApi) ColorerFunc() ColorerFunc {
	return func(ns string, _ Header, re RowEvent) tcell.Color {
		return tcell.ColorCadetBlue
	}
}

// EnvoyAPIRes represents an envoy debug api resource.
type EnvoyApiRes struct {
	Name string
}

// GetObjectKind returns a schema object.
func (EnvoyApiRes) GetObjectKind() schema.ObjectKind {
	return nil
}

// DeepCopyObject returns a container copy.
func (h EnvoyApiRes) DeepCopyObject() runtime.Object {
	return h
}
