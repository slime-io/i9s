package dao

import (
	"context"
	"github.com/derailed/k9s/internal"
	"github.com/derailed/k9s/internal/render"
	"github.com/rs/zerolog/log"
	"k8s.io/apimachinery/pkg/runtime"
)

// IptableInfo
type IptableInfo struct {
	NonResource
}

// List
func (i IptableInfo) List(ctx context.Context, ns string) ([]runtime.Object, error) {
	oo := make([]runtime.Object, 0)
	command := []string{
		"IptableInfo",
	}

	path, ok := ctx.Value(internal.KeyPath).(string)
	if !ok {
		log.Error().Msgf("Expecting a string but got %T", ctx.Value("path"))
		return oo, nil
	}

	for _, f := range command {
		oo = append(oo, render.IptableInfoRes{Name: f, Parent: path})
	}
	return oo, nil
}
