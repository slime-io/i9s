package dao

import (
	"context"
	"github.com/derailed/k9s/internal"
	"github.com/derailed/k9s/internal/render"
	"github.com/rs/zerolog/log"
	"k8s.io/apimachinery/pkg/runtime"
)

// ProxyInfoEx
type ProxyInfo struct {
	NonResource
}

// List
func (i ProxyInfo) List(ctx context.Context, ns string) ([]runtime.Object, error) {
	oo := make([]runtime.Object, 0)
	// "iptables -t nat -S",
	command := []string{
		"istioctl proxy-status",
		"istioctl proxy-config bootstrap",
		"istioctl proxy-config cluster",
		"istioctl proxy-config all",
		"istioctl proxy-config endpoints",
		"istioctl proxy-config listener",
		"istioctl proxy-config log",
		"istioctl proxy-config route",
	}

	path, ok := ctx.Value(internal.KeyPath).(string)
	if !ok {
		log.Error().Msgf("Expecting a string but got %T", ctx.Value("path"))
		return oo, nil
	}

	for _, f := range command {
		oo = append(oo, render.ProxyInfoRes{Name: f, Parent: path})
	}
	return oo, nil
}
