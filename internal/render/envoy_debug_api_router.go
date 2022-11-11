package render

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type EnvoyApiRouter struct{}

func (i EnvoyApiRouter) IsGeneric() bool {
	return false
}

func (i EnvoyApiRouter) Render(o interface{}, ns string, r *Row) error {

	res, ok := o.(EnvoyApiRouterRes)
	if !ok {
		return fmt.Errorf("expected EnvoyAPIRes, but got %T", o)
	}
	r.ID, r.Fields = res.Name, append(r.Fields, res.Name)
	return nil
}

func (i EnvoyApiRouter) Header(ns string) Header {
	return Header{
		HeaderColumn{Name: "api"},
	}
}

func (i EnvoyApiRouter) ColorerFunc() ColorerFunc {
	return func(ns string, _ Header, re RowEvent) tcell.Color {
		return tcell.ColorCadetBlue
	}
}

type EnvoyApiRouterRes struct {
	Name string
}

func (EnvoyApiRouterRes) GetObjectKind() schema.ObjectKind {
	return nil
}

// DeepCopyObject returns a container copy.
func (i EnvoyApiRouterRes) DeepCopyObject() runtime.Object {
	return i
}
