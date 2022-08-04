package view

import (
	"context"
	"github.com/derailed/k9s/internal/client"
	"github.com/derailed/k9s/internal/render"
	"github.com/derailed/k9s/internal/ui"
	"github.com/gdamore/tcell/v2"
)

type IstioXdsPushStatsView struct {
	ResourceViewer
}

func NewIstioXdsPushStatsView(gvr client.GVR) ResourceViewer {
	i := IstioXdsPushStatsView{
		ResourceViewer: NewBrowser(gvr),
	}
	i.GetTable().SetColorerFn(render.IstioProxyID{}.ColorerFunc())
	i.GetTable().SetBorderFocusColor(tcell.ColorMediumSpringGreen)
	i.GetTable().SetSelectedStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorMediumSpringGreen).Attributes(tcell.AttrNone))
	i.SetContextFn(i.chartContext)
	i.GetTable().SetEnterFn(i.enter)
	return &i
}

func (i *IstioXdsPushStatsView) chartContext(ctx context.Context) context.Context {
	return ctx
}

func (i *IstioXdsPushStatsView) enter(app *App, model ui.Tabular, gvr, path string) {

}
