package render

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type IstioXdsPushStats struct{}

func (i IstioXdsPushStats) IsGeneric() bool {
	return false
}

func (i IstioXdsPushStats) Render(o interface{}, ns string, r *Row) error {
	res, ok := o.(IstioXdsPushStatsRes)
	if !ok {
		return fmt.Errorf("expected IstioConfigRes, but got %T", o)
	}
	r.ID = res.Name
	r.Fields = append(r.Fields, res.Name, res.Metric[res.Name])

	//r.ID, r.Fields = res.Name, append(r.Fields, "TIMES")
	//r.Fields = append(r.Fields, res.Metric["cds"])
	//r.Fields = append(r.Fields, res.Metric["eds"])
	//r.Fields = append(r.Fields, res.Metric["lds"])
	//r.Fields = append(r.Fields, res.Metric["rds"])
	return nil
}

func (i IstioXdsPushStats) Header(ns string) Header {

	return Header{
		HeaderColumn{Name: "name"},
		HeaderColumn{Name: "metrics"},
	}
	//return Header{
	//	HeaderColumn{Name: "XDS PUSH"},
	//	HeaderColumn{Name: "CDS"},
	//	HeaderColumn{Name: "EDS"},
	//	HeaderColumn{Name: "LDS"},
	//	HeaderColumn{Name: "RDS"},
	//}
}

func (i IstioXdsPushStats) ColorerFunc() ColorerFunc {
	return func(ns string, _ Header, re RowEvent) tcell.Color {
		return tcell.ColorCadetBlue
	}
}

// IstioConfigRes represents an envoy debug api resource.
type IstioXdsPushStatsRes struct {
	Name     string
	Metric map[string]string
}

// GetObjectKind returns a schema object.
func (IstioXdsPushStatsRes) GetObjectKind() schema.ObjectKind {
	return nil
}

// DeepCopyObject returns a container copy.
func (h IstioXdsPushStatsRes) DeepCopyObject() runtime.Object {
	return h
}
