package render

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"strings"
)

type IstioPilot struct{}

func (i IstioPilot) IsGeneric() bool {
	return false
}

func (i IstioPilot) Render(o interface{}, ns string, r *Row) error {
	res, ok := o.(IstioPilotRes)
	if !ok {
		return fmt.Errorf("expected IstioConfigRes, but got %T", o)
	}

	// 112#sidecarz#istio-system/pilot-x
	id := strings.Join([]string{res.Parent, res.Name}, "#")
	r.ID, r.Fields = id, append(r.Fields, res.Name)
	return nil
}

func (i IstioPilot) Header(ns string) Header {
	return Header{
		HeaderColumn{Name: "pilot instances"},
	}
}

func (i IstioPilot) ColorerFunc() ColorerFunc {
	return func(ns string, _ Header, re RowEvent) tcell.Color {
		return tcell.ColorCadetBlue
	}
}

// IstioConfigRes represents an envoy debug api resource.
type IstioPilotRes struct {
	Name     string
	Parent   string
	Revision string
	Api      string
}

// GetObjectKind returns a schema object.
func (IstioPilotRes) GetObjectKind() schema.ObjectKind {
	return nil
}

// DeepCopyObject returns a container copy.
func (h IstioPilotRes) DeepCopyObject() runtime.Object {
	return h
}
