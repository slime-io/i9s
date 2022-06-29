package view

import (
	"context"
	"encoding/json"
	"github.com/derailed/k9s/internal/client"
	"github.com/derailed/k9s/internal/render"
	"github.com/derailed/k9s/internal/ui"
	"github.com/gdamore/tcell/v2"
	"github.com/rs/zerolog/log"
)

type pilot2poryx struct {
	Name string
	IDs  []string
	Rev  string
	Api  string
}

type IstioPodView struct {
	ResourceViewer
}

func NewIstioPodView(gvr client.GVR) ResourceViewer {
	c := IstioPodView{
		ResourceViewer: NewBrowser(gvr),
	}
	c.GetTable().SetColorerFn(render.EnvoyApi{}.ColorerFunc())
	c.GetTable().SetBorderFocusColor(tcell.ColorMediumSpringGreen)
	c.GetTable().SetSelectedStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorMediumSpringGreen).Attributes(tcell.AttrNone))
	c.SetContextFn(c.chartContext)
	c.GetTable().SetEnterFn(c.enter)
	return &c
}

func (i *IstioPodView) chartContext(ctx context.Context) context.Context {
	log.Debug().Msgf("GetSelectedItem %s in IstioApiView.chartContext", i.GetTable().GetSelectedItem())

	rev, api, instance, err := parseIstioPodView(i.GetTable().GetSelectedItem())
	if err != nil {
		log.Error().Msg("get err %s in enter parse")
		return ctx
	}

	ids := execI9sCmdWithHttp(i, api, instance, rev)
	m := &pilot2poryx{
		Name: instance,
		Rev:  rev,
		Api:  api,
	}
	if len(ids) > 0 {
		for _, id := range ids {
			m.IDs = append(m.IDs, id)
		}
		b, err := json.Marshal(m)
		if err != nil {
			log.Error().Msgf("marshal %+v err %s", m, err)
			return ctx
		}
		return context.WithValue(ctx, "parent", string(b))
	}
	return ctx
}

func (i *IstioPodView) enter(app *App, model ui.Tabular, gvr, path string) {

	rev, api, instance, err := parseIstioPodView(path)
	if err != nil {
		log.Error().Msg("get err %s in IstioPodView enter parse")
		return
	}
	log.Info().Msgf("get rev %s, api %s, instance %s", rev, api, instance)

	api, err = formatIstioAPI(api)
	if err != nil {
		log.Error().Msgf("get err %s in IstioApiView enter", err)
		return
	}
	if needProxyID(api) {
		proxyView := NewIstioProxyIDView(client.NewGVR("proxyID"))
		proxyView.SetContextFn(i.chartContext)
		if err := i.App().inject(proxyView); err != nil {
			i.App().Flash().Err(err)
		}
	} else {
		execi9sCmd(i, api, instance, rev, "")
	}
}

// ---
type AdsClient struct {
	ConnectionID string `json:"connectionId"`
}

type Clients struct {
	Total     int         `json:"totalClients"`
	Connected []AdsClient `json:"clients,omitempty"`
}
