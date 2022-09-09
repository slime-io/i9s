package view

import (
	"context"
	"github.com/derailed/k9s/internal/client"
	"github.com/derailed/k9s/internal/render"
	"github.com/derailed/k9s/internal/ui"
	"github.com/gdamore/tcell/v2"
	"github.com/rs/zerolog/log"
)

const (
	container     = "discovery"
	containerPort = "15014"
)

type IstioApiView struct {
	ResourceViewer
}

func NewIstioApiView(gvr client.GVR) ResourceViewer {
	c := IstioApiView{
		ResourceViewer: NewBrowser(gvr),
	}
	c.GetTable().SetColorerFn(render.EnvoyApi{}.ColorerFunc())
	c.GetTable().SetBorderFocusColor(tcell.ColorMediumSpringGreen)
	c.GetTable().SetSelectedStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorMediumSpringGreen).Attributes(tcell.AttrNone))
	c.GetTable().SetEnterFn(c.enter)
	return &c
}

func (i *IstioApiView) chartContext(ctx context.Context) context.Context {
	// 112#istio/config_dump
	log.Info().Msgf("GetSelectedItem %s in IstioApiView.chartContext", i.GetTable().GetSelectedItem())
	return context.WithValue(ctx, "parent", i.GetTable().GetSelectedItem())
}

func (i *IstioApiView) enter(app *App, model ui.Tabular, gvr, path string) {
	rev, api, err := parseRevWithAPI(path)
	if err != nil {
		log.Error().Msgf("get err in IstioApiView enter, %s", err)
		return
	}
	log.Info().Msgf("get rev %s, api %s in path %s", rev, api, path)

	if api, err = formatIstioAPI(api); err != nil {
		log.Error().Msgf("get err %s in IstioApiView enter", err)
		return
	}

	// need choose one pilot instance
	if needNNK(api) {
		configzView := NewIstioConfigzView(client.NewGVR("configz"))
		configzView.SetContextFn(i.chartContext)
		if err := i.App().inject(configzView); err != nil {
			i.App().Flash().Err(err)
			return
		}
	} else if needPodSelected(api) {
		pilotview := NewIstioPodView(client.NewGVR("pilot"))
		pilotview.SetContextFn(i.chartContext)
		if err := i.App().inject(pilotview); err != nil {
			i.App().Flash().Err(err)
			return
		}
	} else {
		execi9sCmd(i, api, "", rev, "")
	}
}
