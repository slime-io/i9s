package dao

import (
	"context"
	"github.com/derailed/k9s/internal"
	"github.com/derailed/k9s/internal/render"
	"github.com/rs/zerolog/log"
	"k8s.io/apimachinery/pkg/runtime"
)

type IstioApi struct {
	NonResource
}

// List
func (i IstioApi) List(ctx context.Context, ns string) ([]runtime.Object, error) {
	oo := make([]runtime.Object, 0)
	command := []string{
		"istio/adsz",
		"istio/authorizationz",
		"istio/cachez",
		"istio/configz",
		"istio/connections",
		"istio/endpointShardz",
		"istio/endpointz",
		"istio/inject",
		"istio/instancesz",
		"istio/mcsz",
		"istio/mesh",
		"istio/networkz",
		"istio/push_status",
		"istio/pushcontext",
		"istio/registryz",
		"istio/resourcesz",
		"istio/syncz",
		"istio/telemetryz",
		"istio/edsz",
		"istio/sidecarz",
		"istio/config_dump",
		"istio/metrics",
		"istio/xds_push_stats",
		"istio/configzEx",
		"istio/adszEx",
	}
	parent, ok := ctx.Value(internal.Parent).(string)
	if !ok {
		log.Error().Msgf("Expecting a string but got %T", ctx.Value(internal.Parent))
		return oo, nil
	}
	for _, f := range command {
		oo = append(oo, render.IstioApiRes{Name: f, Parent: parent})
	}
	return oo, nil
}
