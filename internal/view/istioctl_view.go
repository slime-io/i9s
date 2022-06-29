package view

import (
	"context"
	"github.com/derailed/k9s/internal/client"
	"github.com/derailed/k9s/internal/render"
	"github.com/derailed/k9s/internal/ui"
	"github.com/gdamore/tcell/v2"
	"github.com/rs/zerolog/log"
)

// IstioView represents a istio command viewer.
type IstioView struct {
	ResourceViewer
}

func NewIstio(gvr client.GVR) ResourceViewer {
	c := IstioView{
		ResourceViewer: NewBrowser(gvr),
	}
	c.GetTable().SetColorerFn(render.Istio{}.ColorerFunc())
	c.GetTable().SetBorderFocusColor(tcell.ColorMediumSpringGreen)
	c.GetTable().SetSelectedStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorMediumSpringGreen).Attributes(tcell.AttrNone))
	c.AddBindKeysFn(c.bindKeys)
	c.SetContextFn(c.chartContext)
	c.GetTable().SetEnterFn(c.showIstioConfig)
	return &c
}

func (i *IstioView) chartContext(ctx context.Context) context.Context {
	rev := i.GetTable().GetSelectedItem()
	return context.WithValue(ctx, "parent", rev)
}

func (i *IstioView) bindKeys(aa ui.KeyActions) {
	aa.Add(ui.KeyActions{
		ui.KeyM: ui.NewKeyAction("isito api", i.debug, true),
	})
}

func (i *IstioView) debug(evt *tcell.EventKey) *tcell.EventKey {
	sel := i.GetTable().GetSelectedItem()
	log.Debug().Msgf("get sel %s in debug", sel)
	if sel == "" {
		return evt
	}
	ida := NewIstioApiView(client.NewGVR("ida"))
	ida.SetContextFn(i.chartContext)
	if err := i.App().inject(ida); err != nil {
		i.App().Flash().Err(err)
		return evt
	}
	return nil
}

func (i *IstioView) showIstioConfig(app *App, model ui.Tabular, gvr, path string) {
	//sel := i.GetTable().GetSelectedItem()
	//log.Info().Msgf("get sel %s in debug in showIstioConfig ", sel)
	ic := NewIstioConfigView(client.NewGVR("ic"))
	ic.SetContextFn(i.chartContext)
	if err := i.App().inject(ic); err != nil {
		i.App().Flash().Err(err)
	}
}
