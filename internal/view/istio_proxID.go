package view

import (
	"context"
	"github.com/derailed/k9s/internal/client"
	"github.com/derailed/k9s/internal/render"
	"github.com/derailed/k9s/internal/ui"
	"github.com/gdamore/tcell/v2"
	"github.com/rs/zerolog/log"
)

type IstioProxyIDView struct {
	ResourceViewer
}

func NewIstioProxyIDView(gvr client.GVR) ResourceViewer {
	c := IstioProxyIDView{
		ResourceViewer: NewBrowser(gvr),
	}
	c.GetTable().SetColorerFn(render.IstioProxyID{}.ColorerFunc())
	c.GetTable().SetBorderFocusColor(tcell.ColorMediumSpringGreen)
	c.GetTable().SetSelectedStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorMediumSpringGreen).Attributes(tcell.AttrNone))
	c.SetContextFn(c.chartContext)
	c.GetTable().SetEnterFn(c.enter)
	return &c
}

func (i *IstioProxyIDView) chartContext(ctx context.Context) context.Context {
	return ctx
}

func (i *IstioProxyIDView) enter(app *App, model ui.Tabular, gvr, path string) {
	// pilot 列表 enter ， 如果是 adsz 则直接执行 cmd, 如果是 sidecarz 需要选择 proxyID
	// istio-system/pilot-1  # httpbin-7476587fbf-zdh9j.apigw-demo-14118

	api, pilot, proxy := parseIstioProxyIDViewID(path)
	log.Info().Msgf("get api %s pilot %s, proxy %s", api, pilot, proxy)

	api, err := formatIstioAPI(api)
	if err != nil {
		log.Error().Msgf("get err %s in IstioApiView enter", err)
		return
	}
	execi9sCmd(i, api, pilot, "", proxy)
}
