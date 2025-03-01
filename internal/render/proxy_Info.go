package render

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"strings"
)

type ProxyInfo struct{}

func (i ProxyInfo) IsGeneric() bool {
	return false
}

func (i ProxyInfo) Render(o interface{}, ns string, r *Row) error {

	res, ok := o.(ProxyInfoRes)
	if !ok {
		return fmt.Errorf("expected ProxyInfoRes, but got %T", o)
	}
	id := strings.Join([]string{res.Parent, res.Name}, "#")
	r.ID, r.Fields = id, append(r.Fields, res.Name)
	return nil
}

func (i ProxyInfo) Header(ns string) Header {
	return Header{
		HeaderColumn{Name: "istioctl command items"},
	}
}

func (i ProxyInfo) ColorerFunc() ColorerFunc {
	return func(ns string, _ Header, re RowEvent) tcell.Color {
		return tcell.ColorCadetBlue
	}
}

type ProxyInfoRes struct {
	Name   string
	Parent string
}

// GetObjectKind returns a schema object.
func (ProxyInfoRes) GetObjectKind() schema.ObjectKind {
	return nil
}

// DeepCopyObject returns a container copy.
func (h ProxyInfoRes) DeepCopyObject() runtime.Object {
	return h
}
