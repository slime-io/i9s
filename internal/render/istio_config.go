package render

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"strings"
)

type IstioConfig struct{}

func (i IstioConfig) IsGeneric() bool {
	return false
}

func (i IstioConfig) Render(o interface{}, ns string, r *Row) error {

	res, ok := o.(IstioConfigRes)
	if !ok {
		return fmt.Errorf("expected IstioConfigRes, but got %T", o)
	}
	id := strings.Join([]string{res.Parent, res.Name}, "#")
	r.ID, r.Fields = id, append(r.Fields, res.Name)
	return nil
}

func (i IstioConfig) Header(ns string) Header {
	return Header{
		HeaderColumn{Name: "Istio Configuration Describe"},
	}
}

func (i IstioConfig) ColorerFunc() ColorerFunc {
	return func(ns string, _ Header, re RowEvent) tcell.Color {
		return tcell.ColorCadetBlue
	}
}

// IstioConfigRes represents an envoy debug api resource.
type IstioConfigRes struct {
	Name   string
	Rev    string
	Parent string
}

// GetObjectKind returns a schema object.
func (IstioConfigRes) GetObjectKind() schema.ObjectKind {
	return nil
}

// DeepCopyObject returns a container copy.
func (h IstioConfigRes) DeepCopyObject() runtime.Object {
	return h
}
