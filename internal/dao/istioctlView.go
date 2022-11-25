package dao

import (
	"context"
	"github.com/derailed/k9s/internal"
	"github.com/derailed/k9s/internal/render"
	"github.com/rs/zerolog/log"
	"k8s.io/apimachinery/pkg/runtime"
)

// IstioctlView
type IstioctlView struct {
	NonResource
}

// List
func (i IstioctlView) List(ctx context.Context, ns string) ([]runtime.Object, error) {
	oo := make([]runtime.Object, 0)
	// "iptables -t nat -S",
	command := []string{
		"istioctl version",
		"istioctl admin log",
		"istioctl profile list",
		"istioctl operator dump",
		"istioctl x revision list",
		"istioctl x config list",
		"istioctl x injector list",
		"istioctl x precheck",
	}

	parent, ok := ctx.Value(internal.Parent).(string)
	if !ok {
		log.Error().Msgf("Expecting a string but got %T", ctx.Value(internal.Parent))
		return oo, nil
	}
	for _, f := range command {
		oo = append(oo, render.IstioctlViewRes{Name: f, Parent: parent})
	}
	return oo, nil
}
