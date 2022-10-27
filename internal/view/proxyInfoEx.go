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

type ProxyinfoExView struct {
	ResourceViewer
}

func NewProxyinfoExView(gvr client.GVR) ResourceViewer {
	c := ProxyinfoExView{
		ResourceViewer: NewBrowser(gvr),
	}
	c.GetTable().SetColorerFn(render.ProxyInfoEx{}.ColorerFunc())
	c.GetTable().SetBorderFocusColor(tcell.ColorMediumSpringGreen)
	c.GetTable().SetSelectedStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorMediumSpringGreen).Attributes(tcell.AttrNone))
	c.SetContextFn(c.chartContext)
	c.GetTable().SetEnterFn(c.enter)
	return &c
}

func (i *ProxyinfoExView) chartContext(ctx context.Context) context.Context {
	return ctx
}

func (i *ProxyinfoExView) enter(app *App, model ui.Tabular, gvr, path string) {

	path, cf, err := parseInfoView(path)
	if err != nil {
		log.Error().Msgf("get error in parseIptableInfoView, %s", err)
		return
	}

	if strings.HasPrefix(cf, "iptables") {
		str := "sudo iptables -t nat -S"
		i.GetIptablesInfo(path, str)
	}
}

func (i *ProxyinfoExView) GetIptablesInfo(path, str string) {

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

	details, ok := pod.Annotations["envoy.io/iptablesDetail"]
	if ok {
		meta.IptablesDetails = details
	}
	parameters := pod.Annotations["envoy.io/iptablesParams"]
	if ok {
		meta.IptablesParams = parameters
	}

	info := i.GetIptablesInContainer(path, str)
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
	IptablesParams  string `json:"iptables_params,omitempty"`
	IptablesRule    string `json:"iptables_rule,omitempty"`
}

func (i *ProxyinfoExView) GetIptablesInContainer(path, str string) string {
	parts := strings.Split(path, "/")
	if len(parts) != 2 {
		return ""
	}
	log.Info().Msgf("get path %s in GetIptablesInContainer", parts)
	info := KubectlProxyCmd(parts[1], parts[0], str)
	return info
}

func KubectlProxyCmd(name, ns, str string) string {

	cmd := fmt.Sprintf("kubectl exec %s -n %s -c istio-proxy -- %s", name, ns, str)
	log.Info().Msgf("prepare exec cmd: %s", cmd)
	out, err := exec.Command("sh", []string{"-c", str}...).Output()
	if err != nil {
		log.Error().Msgf("execProxyCmd err, %s", err)
		return ""
	}
	return string(out)
}

func parseInfoView(s string) (string, string, error) {
	parts := strings.Split(s, "#")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("except 2 items in %s", s)
	}

	return parts[0], parts[1], nil
}
