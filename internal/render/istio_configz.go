package render

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"strings"
)

type IstioConfigz struct{}

func (i IstioConfigz) IsGeneric() bool {
	return false
}

func (i IstioConfigz) Render(o interface{}, ns string, r *Row) error {

	res, ok := o.(IstioConfigzRes)
	if !ok {
		return fmt.Errorf("expected IstioConfigRes, but got %T", o)
	}
	id := strings.Join([]string{res.Kind, res.Namespace, res.Name}, "#")
	id = strings.Join([]string{id, res.Parent}, "#")
	r.ID = id
	r.Fields = Fields{
		res.Kind,
		res.Namespace,
		res.Name,
	}
	return nil
}

func (i IstioConfigz) Header(ns string) Header {
	return Header{
		HeaderColumn{Name: "Kind"},
		HeaderColumn{Name: "Namespace"},
		HeaderColumn{Name: "Name"},
	}
}

func (i IstioConfigz) ColorerFunc() ColorerFunc {
	return func(ns string, _ Header, re RowEvent) tcell.Color {
		return tcell.ColorCadetBlue
	}
}

// IstioConfigRes represents an envoy debug api resource.
type IstioConfigzRes struct {
	Namespace string
	Kind string
	Name string
	Parent string
}

// GetObjectKind returns a schema object.
func (IstioConfigzRes) GetObjectKind() schema.ObjectKind {
	return nil
}

// DeepCopyObject returns a container copy.
func (h IstioConfigzRes) DeepCopyObject() runtime.Object {
	return h
}
