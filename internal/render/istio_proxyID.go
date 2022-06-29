package render

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"strings"
)

type IstioProxyID struct{}

func (i IstioProxyID) IsGeneric() bool {
	return false
}

func (i IstioProxyID) Render(o interface{}, ns string, r *Row) error {
	res, ok := o.(IstioProxyIDRes)
	if !ok {
		return fmt.Errorf("expected IstioConfigRes, but got %T", o)
	}

	// istio-system/pilot-x  ## reviews
	id := strings.Join([]string{res.Revision, res.Api, res.Instance, res.Name}, "#")
	r.ID, r.Fields = id, append(r.Fields, res.Name)
	return nil
}

func (i IstioProxyID) Header(ns string) Header {
	return Header{
		HeaderColumn{Name: "proxy ID"},
	}
}

func (i IstioProxyID) ColorerFunc() ColorerFunc {
	return func(ns string, _ Header, re RowEvent) tcell.Color {
		return tcell.ColorCadetBlue
	}
}

// IstioConfigRes represents an envoy debug api resource.
type IstioProxyIDRes struct {
	Name     string
	Parent   string
	Revision string
	Api      string
	Instance string
}

// GetObjectKind returns a schema object.
func (IstioProxyIDRes) GetObjectKind() schema.ObjectKind {
	return nil
}

// DeepCopyObject returns a container copy.
func (h IstioProxyIDRes) DeepCopyObject() runtime.Object {
	return h
}
