package dao

import (
	"context"
	"github.com/derailed/k9s/internal/render"
	"k8s.io/apimachinery/pkg/runtime"
)

type EnvoyApiRouter struct {
	NonResource
}

func (i EnvoyApiRouter) List(ctx context.Context, ns string) ([]runtime.Object, error) {

	oo := make([]runtime.Object, 0, 2)
	command := []string{
		"envoy/config_dump",
	}
	for _, f := range command {
		oo = append(oo, render.EnvoyApiRouterRes{Name: f})
	}
	return oo, nil
}
