package render

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"strings"
)

type IptableInfo struct{}

func (i IptableInfo) IsGeneric() bool {
	return false
}

func (i IptableInfo) Render(o interface{}, ns string, r *Row) error {

	res, ok := o.(IptableInfoRes)
	if !ok {
		return fmt.Errorf("expected EnvoyAPIRes, but got %T", o)
	}
	id := strings.Join([]string{res.Parent, res.Name}, "#")
	r.ID, r.Fields = id, append(r.Fields, res.Name)
	return nil
}

func (i IptableInfo) Header(ns string) Header {
	return Header{
		HeaderColumn{Name: "iptableInfo"},
	}
}

func (i IptableInfo) ColorerFunc() ColorerFunc {
	return func(ns string, _ Header, re RowEvent) tcell.Color {
		return tcell.ColorCadetBlue
	}
}

// EnvoyAPIRes represents an envoy debug api resource.
type IptableInfoRes struct {
	Name string
	Parent string
}

// GetObjectKind returns a schema object.
func (IptableInfoRes) GetObjectKind() schema.ObjectKind {
	return nil
}

// DeepCopyObject returns a container copy.
func (h IptableInfoRes) DeepCopyObject() runtime.Object {
	return h
}
