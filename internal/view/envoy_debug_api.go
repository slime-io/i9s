package view

import (
	"context"
	"fmt"
	"github.com/derailed/k9s/internal/client"
	"github.com/derailed/k9s/internal/render"
	"github.com/derailed/k9s/internal/ui"
	"github.com/gdamore/tcell/v2"
	"github.com/rs/zerolog/log"
	"strings"
)

type EnvoyApiView struct {
	ResourceViewer
}

func NewEnvoyApiView(gvr client.GVR) ResourceViewer {
	c := EnvoyApiView{
		ResourceViewer: NewBrowser(gvr),
	}
	c.GetTable().SetColorerFn(render.EnvoyApi{}.ColorerFunc())
	c.GetTable().SetBorderFocusColor(tcell.ColorMediumSpringGreen)
	c.GetTable().SetSelectedStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorMediumSpringGreen).Attributes(tcell.AttrNone))
	c.SetContextFn(c.chartContext)
	c.GetTable().SetEnterFn(c.execCmd)
	return &c
}

func (c *EnvoyApiView) chartContext(ctx context.Context) context.Context {
	return ctx
}

func (c *EnvoyApiView) execCmd(app *App, model ui.Tabular, gvr, path string){
	namespace, name := client.Namespaced(c.GetTable().Path)
	cmd := buildIstiocltCmd(path, name, namespace)
	cb := func() {
		opts := shellOpts{
			clear:      false,
			binary:     "sh",
			background: false,
			args:       []string{"-c", cmd},
		}
		if run(c.App(), opts) {
			c.App().Flash().Info("Plugin command launched successfully in envoy view!")
			return
		}
		c.App().Flash().Info("Plugin command failed in envoy view!")
	}
	cb()
}

func buildIstiocltCmd(key, name, namespace string) string {
	parts := strings.Split(key,"/")
	cmd := parts[1]
	if cmd == "config_dump" {
		cmd = "all"
	}
	str := fmt.Sprintf("istioctl proxy-config %s %s.%s -ojson | less", cmd, name, namespace)
	log.Info().Msg(str)
	return str
}