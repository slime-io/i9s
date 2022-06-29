package dao

import (
	"context"
	"github.com/derailed/k9s/internal/render"
	"k8s.io/apimachinery/pkg/runtime"
)

type EnvoyApi struct {
	NonResource
}

// List istioctl proxy-config
func (i EnvoyApi) List(ctx context.Context, ns string) ([]runtime.Object, error) {

	oo := make([]runtime.Object, 0, 2)
	command := []string{
		"envoy/config_dump",
		"envoy/bootstrap",
		"envoy/cluster",
		"envoy/listener",
		"envoy/route",
		"envoy/endpoints",
		"envoy/log",
	}
	for _, f := range command {
		oo = append(oo, render.EnvoyApiRes{Name: f})
	}
	return oo, nil
}