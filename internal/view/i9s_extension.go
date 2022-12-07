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
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"os/exec"
	"strings"
)

type I9sExtensionView struct {
	ResourceViewer
	ParentView string
}

func NewI9sExtensionView(gvr client.GVR, pp string) ResourceViewer {
	c := I9sExtensionView{
		ResourceViewer: NewBrowser(gvr),
		ParentView:     pp,
	}
	c.GetTable().SetColorerFn(render.I9sExtension{}.ColorerFunc())
	c.GetTable().SetBorderFocusColor(tcell.ColorMediumSpringGreen)
	c.GetTable().SetSelectedStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorMediumSpringGreen).Attributes(tcell.AttrNone))
	c.SetContextFn(c.chartContext)
	c.GetTable().SetEnterFn(c.enter)

	return &c
}

func (i *I9sExtensionView) chartContext(ctx context.Context) context.Context {
	return ctx
}

func (i *I9sExtensionView) enter(app *App, model ui.Tabular, gvr, path string) {
	env := i.GetTable().envFn()
	log.Debug().Msgf("get env %+v", env)

	cmd := i.fillEnv(env)
	cmd, err := env.Substitute(cmd)

	if err != nil {
		log.Error().Err(err).Msg("Plugin Args match failed")
		return
	}
	log.Info().Msgf("get cmd in i9s extension view:  %s", cmd)

	out, err := exec.Command("sh", []string{"-c", cmd}...).Output()
	if err != nil {
		log.Error().Msgf("exe cmd err, %s", err)
		return
	}
	view := NewDetails(i.App(), "I9s Extension", "Extension", true).Update(string(out))
	if err := i.App().inject(view); err != nil {
		i.App().Flash().Err(err)
	}
}

func (i *I9sExtensionView) fillEnv(env map[string]string) string {
	// NAME:default#cmd reset --rev $ISTIO_REV
	// NAME:details-v1-7d88846999-lz7ts#kubectl describe pods $NAME -n $NAMESPACE
	// NAME:test#echo $ISTIO_REV

	name := env["NAME"]
	parts := strings.Split(name, "#")
	if len(parts) != 2 {
		return ""
	}
	switch i.ParentView {
	case internal.IstioView:
		if ns := i.getDeploymentNs(parts[0]); ns != "" {
			env["NAMESPACE"] = ns
		}
		env["NAME"] = parts[0]
		env["ISTIO_REV"] = parts[0]
	case "pods":
		env["NAME"] = parts[0]
		env["ISTIO_REV"] = i.getRev(env["NAME"], env["NAMESPACE"])
	default:
		env["NAME"] = parts[0]
	}

	log.Debug().Msgf("env is %+v in fillEnv", env)
	cmd := parts[1]
	return cmd
}

func (i *I9sExtensionView) getDeploymentNs(rev string) string {

	dial, err := i.App().factory.Client().Dial()
	if err != nil {
		log.Error().Msgf("get client err in getDeploymentNs, %s", err)
		return ""
	}
	dps, err := dial.AppsV1().Deployments("").
		List(context.TODO(), metav1.ListOptions{
			LabelSelector: toSelector(map[string]string{"app": "istiod", "istio.io/rev": rev}),
		})
	if err != nil {
		log.Error().Msgf("list deployment err, %s", err)
		return ""
	}

	if len(dps.Items) > 0 {
		return dps.Items[0].Namespace
	}
	return ""
}

func (i *I9sExtensionView) getRev(name, ns string) string {

	po, err := i.App().factory.Get("v1/pods", fmt.Sprintf("%s/%s", ns, name), true, labels.Everything())
	if err != nil {
		log.Error().Msgf("get pods err, %s", err.Error())
		return ""
	}

	pod := v1.Pod{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(po.(*unstructured.Unstructured).Object, &pod)
	if err != nil {
		log.Error().Msgf("Unstructured convert to pods err, %s", err.Error())
		return ""
	} else {
		if !hasProxy(pod) {
			log.Warn().Msgf("pods has no proxy, skip")
			return ""
		}
		if rev, ok := pod.Labels[istioRev]; ok {
			return rev
		}
	}

	objectNs, err := i.App().factory.Get("v1/namespaces", fmt.Sprintf("%s/%s", "-", ns), true, labels.Everything())
	if err != nil {
		log.Error().Msgf("get ns err, %s", err.Error())
		return ""
	}

	namespace := v1.Namespace{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(objectNs.(*unstructured.Unstructured).Object, &namespace)
	if err != nil {
		log.Error().Msgf("Unstructured convert to namespace err, %s", err.Error())
		return ""
	} else {
		if rev, ok := namespace.Labels[istioRev]; ok {
			return rev
		}
	}

	return ""
}

func hasProxy(pod v1.Pod) bool {
	for _, co := range pod.Status.ContainerStatuses {
		if co.Name == "istio-proxy" || co.Name == "gateway-proxy" {
			return true
		}
	}
	return false
}
