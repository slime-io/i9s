package view

import (
	"context"
	"github.com/derailed/k9s/internal/client"
	"github.com/derailed/k9s/internal/render"
	"github.com/derailed/k9s/internal/ui"
	"github.com/gdamore/tcell/v2"
	"github.com/rs/zerolog/log"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
)

type IstioConfigView struct {
	ResourceViewer
}

func NewIstioConfigView(gvr client.GVR) ResourceViewer {
	c := IstioConfigView{
		ResourceViewer: NewBrowser(gvr),
	}
	c.GetTable().SetColorerFn(render.EnvoyApi{}.ColorerFunc())
	c.GetTable().SetBorderFocusColor(tcell.ColorMediumSpringGreen)
	c.GetTable().SetSelectedStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorMediumSpringGreen).Attributes(tcell.AttrNone))
	c.SetContextFn(c.chartContext)
	c.GetTable().SetEnterFn(c.enter)
	return &c
}

func (i *IstioConfigView) chartContext(ctx context.Context) context.Context {
	return ctx
}

func (i *IstioConfigView) enter(app *App, model ui.Tabular, gvr, path string) {

	rev, path, err := parseIstioConfig(path)
	if err != nil {
		log.Error().Msgf("get err in IstioConfigView execCmd, %s", err)
		return
	}
	// configmap inject
	log.Info().Msgf("rev %s path %s is selected in IstioConfigView", rev, path)
	dial, err := i.App().factory.Client().Dial()
	if err != nil {
		log.Error().Msgf("get client err in IstioConfigView, %s", err)
		return
	}
	switch path {
	case "Istiod Configuration":
		i.GetConfigMap(dial, rev)
	case "Injector Mutatingwebhookconfigurations":
		i.GetInject(dial, rev)
	case "Istiod Deployment Manifest":
		i.GetDeployment(dial, rev)
	default:
	}
}

func (i *IstioConfigView) GetConfigMap(dial kubernetes.Interface, rev string) {

	label := map[string]string{istioRev: rev}
	oo, err := i.App().factory.List("v1/configmaps", "", true, labels.Set(label).AsSelector())
	if err != nil {
		log.Error().Msgf("list configmap err, %s", err.Error())
		return
	}

	// find the first cm with prefix "istio"
	item := corev1.ConfigMap{}
	for _, o := range oo {
		var cm corev1.ConfigMap
		err = runtime.DefaultUnstructuredConverter.FromUnstructured(o.(*unstructured.Unstructured).Object, &cm)
		if err != nil {
			log.Error().Msgf("Unstructured convert to configmap err, %s", err.Error())
			return
		}
		if cm.Name == "istio-" + rev || cm.Name == "istio" {
			item = cm
			break
		}
	}
	if item.Name != "" {
		s, err := convert2String(item)
		if err != nil {
			log.Error().Msgf("marshal istio cofiguration configmap to yaml err, %s", err)
			return
		}
		details := NewDetails(i.App(), "YAML", "istio cofiguration configmap", true).Update(s)
		if err := i.App().inject(details); err != nil {
			i.App().Flash().Err(err)
		}
	}
}

func (i *IstioConfigView) GetInject(dial kubernetes.Interface, rev string) {

	label := map[string]string{istioRev: rev}
	selector := labels.Set(label).AsSelector()
	list, err := dial.AdmissionregistrationV1beta1().MutatingWebhookConfigurations().List(context.TODO(), metav1.ListOptions{LabelSelector: selector.String()})
	if err != nil {
		log.Error().Msgf("list mutatingwebhookconfigurations with %s err %s", selector.String(), err)
		return
	}
	if len(list.Items) > 0 {
		s, err := convert2String(list.Items[0])
		if err != nil {
			log.Error().Msgf("marshal istio injector to yaml err, %s", err)
			return
		}
		details := NewDetails(i.App(), "YAML", "isito injector", true).Update(s)
		if err := i.App().inject(details); err != nil {
			i.App().Flash().Err(err)
		}
	}
}

func (i *IstioConfigView) GetDeployment(dial kubernetes.Interface, rev string) {
	label := map[string]string{"istio.io/rev": rev, "app": "istiod"}
	selector := labels.Set(label).AsSelector()
	list, err := dial.AppsV1().Deployments("").List(context.TODO(), metav1.ListOptions{LabelSelector: selector.String()})
	if err != nil {
		log.Error().Msgf("get istiod %s", selector.String())
		return
	}
	if len(list.Items) > 0 {
		s, err := convert2String(list.Items[0])
		if err != nil {
			log.Error().Msgf("marshal istio deployment to yaml err, %s", err)
			return
		}
		details := NewDetails(i.App(), "YAML", "istio deployment", true).Update(s)
		if err := i.App().inject(details); err != nil {
			i.App().Flash().Err(err)
		}
	}
}
