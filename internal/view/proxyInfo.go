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

type ProxyinfoView struct {
	ResourceViewer
}

func NewProxyInfoView(gvr client.GVR) ResourceViewer {
	c := ProxyinfoView{
		ResourceViewer: NewBrowser(gvr),
	}
	c.GetTable().SetColorerFn(render.ProxyInfo{}.ColorerFunc())
	c.GetTable().SetBorderFocusColor(tcell.ColorMediumSpringGreen)
	c.GetTable().SetSelectedStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorMediumSpringGreen).Attributes(tcell.AttrNone))
	c.SetContextFn(c.chartContext)
	c.GetTable().SetEnterFn(c.enter)
	return &c
}

func (i *ProxyinfoView) chartContext(ctx context.Context) context.Context {
	return ctx
}

func (i *ProxyinfoView) enter(app *App, model ui.Tabular, gvr, path string) {

	path, cf, err := parseInfoView(path)
	if err != nil {
		log.Error().Msgf("get error in parseIptableInfoView, %s", err)
		return
	}
	if strings.HasPrefix(cf, "istioctl") {
		i.istioctlProxyCmd(cf, path)
	}
}

func (c *ProxyinfoView) istioctlProxyCmd(cf, path string) {

	namespace, name := client.Namespaced(path)
	cmd := fmt.Sprintf("%s %s.%s | less", cf, name, namespace)

	log.Info().Msgf("prepare exec cmd: %s", cmd)
	cb := func() {
		opts := shellOpts{
			clear:      false,
			binary:     "sh",
			background: false,
			args:       []string{"-c", cmd},
		}
		if run(c.App(), opts) {
			c.App().Flash().Info("command launched successfully in proxyInfo!")
			return
		}
		c.App().Flash().Info("command failed in proxyInfo!")
	}
	cb()
}
