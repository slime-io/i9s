package dao

import (
	"context"
	"github.com/derailed/k9s/internal"
	"github.com/derailed/k9s/internal/config"
	"github.com/derailed/k9s/internal/render"
	"github.com/rs/zerolog/log"
	"k8s.io/apimachinery/pkg/runtime"
	"strings"
)

type I9sExtension struct {
	NonResource
}

// List
func (i *I9sExtension) List(ctx context.Context, ns string) ([]runtime.Object, error) {
	oo := make([]runtime.Object, 0)

	typ, ok := ctx.Value(internal.ExtensionType).(string)
	if !ok {
		log.Error().Msgf("want to get stirng extensionType, but get %T", ctx.Value(internal.ExtensionType))
		return oo, nil
	}
	log.Info().Msgf("get typ %s", typ)
	command := getCommand(string(typ))
	log.Info().Msgf("get command %+v", command)
	parent, ok := ctx.Value(internal.IstioRev).(string)
	if !ok {
		log.Error().Msgf("want to get stirng extensionType, but get %T", ctx.Value(internal.IstioRev))
		return oo, nil
	}

	for _, f := range command {
		oo = append(oo, render.I9sExtensionRes{Name: f, Parent: string(parent)})
	}
	return oo, nil
}

func getCommand(typ string) []string {
	var cmds []string
	arr := config.GetViewToExtension()[typ]
	for _, ext := range arr {
		cmds = append(cmds, strings.Join(ext.Args, " "))
	}
	return cmds
}
