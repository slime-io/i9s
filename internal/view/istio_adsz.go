package view

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/derailed/k9s/internal/client"
	"github.com/derailed/k9s/internal/render"
	"github.com/derailed/k9s/internal/ui"
	"github.com/gdamore/tcell/v2"
	"github.com/rs/zerolog/log"
	"os/exec"
	"strings"
)

type IstioAdszView struct {
	ResourceViewer
}

func NewIstioAdszView(gvr client.GVR) ResourceViewer {
	c := IstioAdszView{
		ResourceViewer: NewBrowser(gvr),
	}
	c.GetTable().SetColorerFn(render.IstioAdsz{}.ColorerFunc())
	c.GetTable().SetBorderFocusColor(tcell.ColorMediumSpringGreen)
	c.GetTable().SetSelectedStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorMediumSpringGreen).Attributes(tcell.AttrNone))
	c.SetContextFn(c.chartContext)
	c.GetTable().SetEnterFn(c.enter)
	return &c
}

func (i *IstioAdszView) chartContext(ctx context.Context) context.Context {
	return ctx
}

func (i *IstioAdszView) enter(app *App, model ui.Tabular, gvr, path string) {
	namespace, name, connectionId, err := parseAdszPath(path)
	if err != nil {
		log.Error().Msgf("get err in parseAdszCtx %s", err)
		return
	}

	b := getClientInfo(namespace, name, connectionId)
	details := NewDetails(i.App(), "json", "info", true)
	details.Update(string(b))
	if err := i.App().inject(details); err != nil {
		i.App().Flash().Err(err)
	}
}

func getClientInfo(namespace, name, connectionId string ) []byte {
	var b []byte
	config := getAdszConfig(namespace, name)
	if config == "" {
		return b
	}
	adsz := parseAdszConfig(config)
	if adsz.TotalClients == 0 {
		return b
	}
	clientInfo := parseAdszClient(connectionId, adsz)
	b, err := json.MarshalIndent(clientInfo, "", "  ")
	if err != nil {
		log.Error().Msgf("marshal clientInfo err %s", err.Error())
		return []byte{}
	}
	return b
}


func parseAdszPath(path string) (string, string, string, error) {
	parts := strings.Split(path, "#")
	if len(parts) != 2 {
		return "", "", "", fmt.Errorf("invalid path %s", path)
	}
	nn := strings.Split(parts[0], "/")
	if len(nn) != 2 {
		return "", "", "", fmt.Errorf("invalid path %s", path)
	}
	return nn[0], nn[1], parts[1], nil
}


type Adsz struct {
	TotalClients int32 	`json:"totalClients"`
	Clients []Client `json:"clients"`
}

type Client struct {
	ConnectionId string `json:"connectionId"`
	ConnectedAt string `json:"connectedAt"`
	Address string `json:"address"`
	Metadata interface{} `json:"metadata"`
	Watches interface{} `json:"watches"`
}


func getAdszConfig(namespace, name string) string {
	str := fmt.Sprintf("kubectl exec %s -n %s -- curl localhost:15014/debug/adsz -s", name, namespace)
	log.Debug().Msgf("exec %s", str)
	out, err := exec.Command("sh", []string{"-c", str}...).Output()
	if err != nil {
		log.Error().Msgf("exe cmd err, %s", err)
		return ""
	}
	return string(out)
}

func parseAdszConfig(c string) Adsz {
	adsz := Adsz{}
	if err := json.Unmarshal([]byte(c), &adsz); err != nil {
		log.Error().Msgf("unmarshal adsz to struct %d", err.Error())
		return Adsz{}
	}
	return adsz
}

func parseAdszClient(connectionId string, adsz Adsz) Client {
	for _, client := range adsz.Clients {
		if client.ConnectionId == connectionId {
			return client
		}
	}
	return  Client{}
}


