package view

import (
	"github.com/derailed/k9s/internal/client"
	"github.com/gdamore/tcell/v2"
	"github.com/rs/zerolog/log"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
)

func (p *Pod) showProxyConfig(evt *tcell.EventKey) *tcell.EventKey {
	path := p.GetTable().GetSelectedItem()
	if path == "" {
		return evt
	}
	po, err := p.App().factory.Get("v1/pods", p.GetTable().GetSelectedItem(), true, labels.Everything())
	if err != nil {
		log.Error().Msgf("get pods %s err, %s", p.GetTable().GetSelectedItem(), err.Error())
		return evt
	}
	pod := v1.Pod{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(po.(*unstructured.Unstructured).Object, &pod)
	if err != nil {
		log.Error().Msgf("Unstructured convert to configmap err, %s", err.Error())
		return evt
	}
	has := false
	for _, co := range pod.Status.ContainerStatuses {
		if co.Name == "istio-proxy" {
			has = true
			break
		}
	}
	if ! has {
		log.Debug().Msgf("pod %s has no istio-proxy, skip", p.GetTable().GetSelectedItem())
		return evt
	}

	log.Info().Msgf("get path %s in showProxyConfig", path)
	eda := NewEnvoyApiView(client.NewGVR("eda"))
	eda.SetContextFn(p.coContext)
	if err := p.App().inject(eda); err != nil {
		p.App().Flash().Err(err)
		return evt
	}
	return nil
}
