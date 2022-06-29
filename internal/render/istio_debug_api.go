package render

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"strings"
)

type IstioApi struct{}

func (i IstioApi) IsGeneric() bool {
	return false
}

func (i IstioApi) Render(o interface{}, ns string, r *Row) error {

	res, ok := o.(IstioApiRes)
	if !ok {
		return fmt.Errorf("expected EnvoyAPIRes, but got %T", o)
	}
	id := strings.Join([]string{res.Parent, res.Name}, "#")
	r.ID, r.Fields = id, append(r.Fields, res.Name)
	return nil
}

func (i IstioApi) Header(ns string) Header {
	return Header{
		HeaderColumn{Name: "api"},
	}
}

func (i IstioApi) ColorerFunc() ColorerFunc {
	return func(ns string, _ Header, re RowEvent) tcell.Color {
		return tcell.ColorCadetBlue
	}
}

// IstioApiRes represents an istio debug api resource.
type IstioApiRes struct {
	Name   string
	Rev    string
	Parent string
}

// GetObjectKind returns a schema object.
func (IstioApiRes) GetObjectKind() schema.ObjectKind {
	return nil
}

// DeepCopyObject returns a container copy.
func (h IstioApiRes) DeepCopyObject() runtime.Object {
	return h
}
