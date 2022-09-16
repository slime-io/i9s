package render

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"strings"
)

type IstioAdsz struct{}

func (i IstioAdsz) IsGeneric() bool {
	return false
}

func (i IstioAdsz) Render(o interface{}, ns string, r *Row) error {

	res, ok := o.(IstioAdszRes)
	if !ok {
		return fmt.Errorf("expected IstioAdszRes, but got %T", o)
	}
	id := strings.Join([]string{res.IstioInstance, res.ConnectionId}, "#")
	r.ID = id
	r.Fields = Fields{
		res.IstioInstance,
		res.ConnectionId,
		res.ConnectedAt,
		res.Address,
	}
	return nil
}

func (i IstioAdsz) Header(ns string) Header {
	return Header{
		HeaderColumn{Name: "XdsServer"},
		HeaderColumn{Name: "ConnectionId"},
		HeaderColumn{Name: "ConnectedAt"},
		HeaderColumn{Name: "Address"},
	}
}

func (i IstioAdsz) ColorerFunc() ColorerFunc {
	return func(ns string, _ Header, re RowEvent) tcell.Color {
		return tcell.ColorCadetBlue
	}
}

// IstioConfigRes represents an envoy debug api resource.
type IstioAdszRes struct {
	ConnectionId string
	ConnectedAt string
	Address string
	IstioInstance string
	Parent string
}

// GetObjectKind returns a schema object.
func (IstioAdszRes) GetObjectKind() schema.ObjectKind {
	return nil
}

// DeepCopyObject returns a container copy.
func (h IstioAdszRes) DeepCopyObject() runtime.Object {
	return h
}
