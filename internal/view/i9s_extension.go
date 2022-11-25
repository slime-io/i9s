package view

import (
	"context"
	"github.com/derailed/k9s/internal/client"
	"github.com/derailed/k9s/internal/render"
	"github.com/derailed/k9s/internal/ui"
	"github.com/gdamore/tcell/v2"
	"github.com/rs/zerolog/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os/exec"
	"strings"
)

type I9sExtensionView struct {
	ResourceViewer
}

func NewI9sExtensionView(gvr client.GVR) ResourceViewer {
	c := I9sExtensionView{
		ResourceViewer: NewBrowser(gvr),
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
	// NAME:default#cmd reset --rev ${ISTIO_REV}
	// NAME:details-v1-7d88846999-lz7ts#kubectl describe pods ${NAME} -n ${NAMESPACE}
	name := env["NAME"]
	parts := strings.Split(name, "#")
	if len(parts) != 2 {
		return ""
	}
	env["ISTIO_REV"] = parts[0]
	env["NAME"] = parts[0]

	if env["NAMESPACE"] == "" {
		if ns := i.getDeploymentNs(parts[0]); ns != "" {
			env["NAMESPACE"] = ns

		}
	}
	return parts[1]
}

func (i *I9sExtensionView) getDeploymentNs(rev string) string {

	dial, err := i.App().factory.Client().Dial()
	if err != nil {
		log.Error().Msgf("get client err in view i9s_extension, %s", err)
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
