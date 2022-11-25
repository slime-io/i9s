package view

import (
	"context"
	"fmt"
	"github.com/derailed/k9s/internal"
	"github.com/derailed/k9s/internal/client"
	"github.com/derailed/k9s/internal/render"
	"github.com/derailed/k9s/internal/ui"
	"github.com/gdamore/tcell/v2"
	"github.com/rs/zerolog/log"
)

type IstioctlView struct {
	ResourceViewer
}

func NewIstioctlView(gvr client.GVR) ResourceViewer {
	c := IstioctlView{
		ResourceViewer: NewBrowser(gvr),
	}
	c.GetTable().SetColorerFn(render.IstioctlView{}.ColorerFunc())
	c.GetTable().SetBorderFocusColor(tcell.ColorMediumSpringGreen)
	c.GetTable().SetSelectedStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorMediumSpringGreen).Attributes(tcell.AttrNone))
	c.SetContextFn(c.chartContext)
	c.GetTable().SetEnterFn(c.enter)
	return &c
}

func (i *IstioctlView) chartContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, internal.Parent, i.GetTable().GetSelectedItem())
}

func (i *IstioctlView) enter(app *App, model ui.Tabular, gvr, path string) {
	_, cmd, err := parseRevWithAPI(path)
	if err != nil {
		log.Error().Msgf("get error in parseRevWithAPI, %s", err)
		return
	}
	i.istioctlCmd(cmd)
}

func (i *IstioctlView) istioctlCmd(item string) {
	cmd := fmt.Sprintf("%s | less", item)
	log.Info().Msgf("prepare exec cmd: %s", cmd)
	cb := func() {
		opts := shellOpts{
			clear:      false,
			binary:     "sh",
			background: false,
			args:       []string{"-c", cmd},
		}
		if run(i.App(), opts) {
			i.App().Flash().Info("command launched successfully in proxyInfo!")
			return
		}
		i.App().Flash().Info("command failed in proxyInfo!")
	}
	cb()
}
