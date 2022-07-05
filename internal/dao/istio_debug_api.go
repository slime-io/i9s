package dao

import (
	"context"
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
	}
	parent, ok := ctx.Value("parent").(string)
	if !ok {
		log.Error().Msgf("Expecting a string but got %T", ctx.Value("parent"))
		return oo, nil
	}
	for _, f := range command {
		oo = append(oo, render.IstioApiRes{Name: f, Parent: parent})
	}
	return oo, nil
}
