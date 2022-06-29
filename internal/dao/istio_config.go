package dao

import (
	"context"
	"github.com/derailed/k9s/internal/render"
	"github.com/rs/zerolog/log"
	"k8s.io/apimachinery/pkg/runtime"
)

// IstioConfig IC
type IstioConfig struct {
	NonResource
}

// List
func (i IstioConfig) List(ctx context.Context, ns string) ([]runtime.Object, error) {
	oo := make([]runtime.Object, 0)
	command := []string{
		"Istiod Configuration",
		"Injector Mutatingwebhookconfigurations",
		"Istiod Deployment Manifest",
	}

	parent, ok := ctx.Value("parent").(string)
	if !ok {
		log.Error().Msgf("Expecting a string but got %T", ctx.Value("parent"))
		return oo, nil
	}

	for _, f := range command {
		oo = append(oo, render.IstioConfigRes{Name: f, Parent: parent})
	}
	return oo, nil
}
