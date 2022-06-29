package render

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// Istio renders istio command to screen.
type Istio struct{}

func (i Istio) IsGeneric() bool {
	return false
}

func (i Istio) Render(o interface{}, ns string, r *Row) error {

	res, ok := o.(IstioRes)
	if !ok {
		return fmt.Errorf("expected *IstioRes, but got %T", o)
	}
	r.ID, r.Fields = res.Name, append(r.Fields, res.Name)
	return nil
}

func (i Istio) Header(ns string) Header {
	return Header{
		HeaderColumn{Name: "REVISION"},
	}
}

func (i Istio) ColorerFunc() ColorerFunc {
	return func(ns string, _ Header, re RowEvent) tcell.Color {
		return tcell.ColorCadetBlue
	}
}

type IstioRes struct {
	Name string
}

// GetObjectKind returns a schema object.
func (IstioRes) GetObjectKind() schema.ObjectKind {
	return nil
}

// DeepCopyObject returns a container copy.
func (h IstioRes) DeepCopyObject() runtime.Object {
	return h
}
