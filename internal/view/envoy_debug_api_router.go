package view

import (
	"context"
	"fmt"
	"github.com/derailed/k9s/internal/client"
	"github.com/derailed/k9s/internal/render"
	"github.com/derailed/k9s/internal/ui"
	"github.com/gdamore/tcell/v2"
	"strings"
)

var usefx_router = true

type EnvoyApiRouterView struct {
	ResourceViewer
}

func NewEnvoyApiRouterView(gvr client.GVR) ResourceViewer {
	c := EnvoyApiRouterView{
		ResourceViewer: NewBrowser(gvr),
	}
	c.GetTable().SetColorerFn(render.EnvoyApiRouter{}.ColorerFunc())
	c.GetTable().SetBorderFocusColor(tcell.ColorMediumSpringGreen)
	c.GetTable().SetSelectedStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorMediumSpringGreen).Attributes(tcell.AttrNone))
	c.SetContextFn(c.chartContext)
	c.AddBindKeysFn(c.bindKeys)
	c.GetTable().SetEnterFn(c.execCmd)
	return &c
}

func (e *EnvoyApiRouterView) bindKeys(aa ui.KeyActions) {
	aa.Add(ui.KeyActions{
		ui.KeyP: ui.NewKeyAction("switch to fx/less", e.switchon, true),
	})
}

func (e *EnvoyApiRouterView) switchon(evt *tcell.EventKey) *tcell.EventKey {
	usefx_router = !usefx_router
	return nil
}

func (e *EnvoyApiRouterView) chartContext(ctx context.Context) context.Context {
	return ctx
}

func (e *EnvoyApiRouterView) execCmd(app *App, model ui.Tabular, gvr, path string) {
	namespace, name := client.Namespaced(e.GetTable().Path)
	cmd := buildKubectlCmd(path, name, namespace)
	cb := func() {
		opts := shellOpts{
			clear:      false,
			binary:     "sh",
			background: false,
			args:       []string{"-c", cmd},
		}
		if run(e.App(), opts) {
			e.App().Flash().Info("Plugin command launched successfully in envoy view!")
			return
		}
		e.App().Flash().Info("Plugin command failed in envoy view!")
	}
	cb()
}

func buildKubectlCmd(key, name, namespace string) string {
	parts := strings.Split(key, "/")
	cmd := parts[1]
	if cmd == "config_dump" {
		str := fmt.Sprintf("kubectl exec %s -n %s -- curl 127.0.0.1:19000/config_dump -s", name, namespace)
		if usefx_router {
			str = str + " | fx"
		} else {
			str = str + " | less"
		}
		return str
	}
	return ""
}
