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

	if !p.containProxy() {
		log.Debug().Msgf("pod %s has no istio-proxy, skip", path)
		return nil
	}

	//log.Info().Msgf("get path %s in showProxyConfig", path)
	eda := NewEnvoyApiView(client.NewGVR("eda"))
	eda.SetContextFn(p.coContext)
	if err := p.App().inject(eda); err != nil {
		p.App().Flash().Err(err)
		return evt
	}
	return nil
}

func (p *Pod) showProxyInfo(evt *tcell.EventKey) *tcell.EventKey {

	path := p.GetTable().GetSelectedItem()
	//log.Info().Msgf("get pods %s in showProxyMetaInfo", path)
	if path == "" {
		return evt
	}

	if !p.containProxy() {
		log.Debug().Msgf("pod %s has no istio-proxy, skip", path)
		return nil
	}

	info := NewProxyInfoView(client.NewGVR("proxyInfo"))
	info.SetContextFn(p.coContext)
	if err := p.App().inject(info); err != nil {
		p.App().Flash().Err(err)
		return evt
	}
	return nil
}

func (p *Pod) showProxyInfoEx(evt *tcell.EventKey) *tcell.EventKey {

	path := p.GetTable().GetSelectedItem()
	//log.Info().Msgf("get pods %s in showProxyMetaInfo", path)
	if path == "" {
		return evt
	}

	if !p.containProxy() {
		log.Debug().Msgf("pod %s has no istio-proxy, skip", path)
		return nil
	}

	info := NewProxyinfoExView(client.NewGVR("proxyInfoEx"))
	info.SetContextFn(p.coContext)
	if err := p.App().inject(info); err != nil {
		p.App().Flash().Err(err)
		return evt
	}
	return nil
}

func (p *Pod) containProxy() bool {
	path := p.GetTable().GetSelectedItem()
	if path == "" {
		return false
	}
	po, err := p.App().factory.Get("v1/pods", p.GetTable().GetSelectedItem(), true, labels.Everything())
	if err != nil {
		log.Error().Msgf("get pods %s err, %s", p.GetTable().GetSelectedItem(), err.Error())
		return false
	}
	pod := v1.Pod{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(po.(*unstructured.Unstructured).Object, &pod)
	if err != nil {
		log.Error().Msgf("Unstructured convert to configmap err, %s", err.Error())
		return false
	}
	has := false
	for _, co := range pod.Status.ContainerStatuses {
		if co.Name == "istio-proxy" {
			has = true
			break
		}
	}
	return has
}
