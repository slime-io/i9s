package dao

import (
	"context"
	"encoding/json"
	"github.com/derailed/k9s/internal"
	"github.com/derailed/k9s/internal/render"
	"github.com/rs/zerolog/log"
	"k8s.io/apimachinery/pkg/runtime"
)

// IstioConfig IC
type IstioProxyID struct {
	NonResource
}

type pilot2poryx struct {
	Name string
	IDs  []string
	Rev  string
	Api  string
}

// List
func (i IstioProxyID) List(ctx context.Context, ns string) ([]runtime.Object, error) {
	oo := make([]runtime.Object, 0)
	// 用conetxt 传递 jso
	// 112#istio/sidecarz#istio-system/pilot-1 => value
	// json
	parent, ok := ctx.Value(internal.Parent).(string)
	if !ok {
		log.Error().Msgf("Expecting a string but got %T", ctx.Value(internal.Parent))
		return oo, nil
	}
	//log.Info().Msgf("get parent %s in dao IstioProxyID", parent)
	m := new(pilot2poryx)
	err := json.Unmarshal([]byte(parent), &m)
	if err != nil {
		log.Error().Msgf("unmarshal from parent %s get err, %s", parent, err)
		return oo, nil
	}

	for _, id := range m.IDs {
		oo = append(oo, render.IstioProxyIDRes{Name: id, Instance: m.Name, Api: m.Api, Revision: m.Rev})
	}
	return oo, nil
}
