package view

import (
	"context"
	"fmt"
	"github.com/derailed/k9s/internal/client"
	"github.com/derailed/k9s/internal/render"
	"github.com/derailed/k9s/internal/ui"
	"github.com/gdamore/tcell/v2"
	"github.com/rs/zerolog/log"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"os/exec"
	"strings"
)

type EnvoyMetaInfoView struct {
	ResourceViewer
}

func NewIptableInfoView(gvr client.GVR) ResourceViewer {
	c := EnvoyMetaInfoView{
		ResourceViewer: NewBrowser(gvr),
	}
	c.GetTable().SetColorerFn(render.IstioProxyID{}.ColorerFunc())
	c.GetTable().SetBorderFocusColor(tcell.ColorMediumSpringGreen)
	c.GetTable().SetSelectedStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorMediumSpringGreen).Attributes(tcell.AttrNone))
	c.SetContextFn(c.chartContext)
	c.GetTable().SetEnterFn(c.enter)
	return &c
}

func (i *EnvoyMetaInfoView) chartContext(ctx context.Context) context.Context {
	return ctx
}

func (i *EnvoyMetaInfoView) enter(app *App, model ui.Tabular, gvr, path string) {

	path, cf, err := parseIptableInfoView(path)
	if err != nil {
		log.Error().Msgf("get error in parseIptableInfoView, %s", err)
		return
	}

	switch cf {
	case "IptableInfo":
		i.GetIptablesInfo(path)
	default:
	}
}

func (i *EnvoyMetaInfoView) GetIptablesInfo(path string) {

	po, err := i.App().factory.Get("v1/pods", path, true, labels.Everything())
	if err != nil {
		log.Error().Msgf("get pods %s err, %s", path, err.Error())
		return
	}
	pod := corev1.Pod{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(po.(*unstructured.Unstructured).Object, &pod)
	if err != nil {
		log.Error().Msgf("Unstructured convert to configmap err, %s", err.Error())
		return
	}
	meta := IptablesMeta{}

	details , ok := pod.Annotations["envoy.io/iptablesDetail"]
	if ok {
		meta.IptablesDetails = details
	}
	parameters := pod.Annotations["envoy.io/iptablesParams"]
	if ok {
		meta.IptablesParams = parameters
	}

	info := i.GetIptablesInContainer(path)
	meta.IptablesRule = info

	content, err := convert2String(meta)
	if err != nil {
		log.Error().Msgf("get err in GetTrafficRules, %s", err)
		return
	}

	view := NewDetails(i.App(), "Iptables Info", "Iptables Info", true).Update(content)
	if err := i.App().inject(view); err != nil {
		i.App().Flash().Err(err)
	}
}

type IptablesMeta struct {
	IptablesDetails string `json:"iptables_details,omitempty"`
	IptablesParams string `json:"iptables_params,omitempty"`
	IptablesRule string `json:"iptables_rule,omitempty"`
}

func (i *EnvoyMetaInfoView) GetIptablesInContainer(path string) string{
	parts := strings.Split(path, "/")
	if len(parts) != 2 {
		return ""
	}
	log.Info().Msgf("get path %s in GetIptablesInContainer", parts)
	info := execProxyCmd(parts[1], parts[0])
	return info
}


func parseIptableInfoView(s string) (string, string, error) {
	parts := strings.Split(s, "#")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("except 2 items in %s", s)
	}

	return parts[0], parts[1], nil
}

func execProxyCmd(name,ns string) string {

	str := fmt.Sprintf("kubectl exec %s -n %s -c istio-proxy -- sudo iptables -t nat -S", name, ns)
	log.Info().Msgf("exec %s", str)
	out, err := exec.Command("sh", []string{"-c", str}...).Output()
	if err != nil {
		log.Error().Msgf("execProxyCmd err, %s", err)
		return ""
	}
	return string(out)
}
