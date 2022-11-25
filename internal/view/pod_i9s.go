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

	has, sidecar := p.containProxy()
	if !has {
		log.Debug().Msgf("pod %s has no proxy, skip", path)
		return nil
	}

	if sidecar {
		eda := NewEnvoyApiView(client.NewGVR("eda"))
		eda.SetContextFn(p.coContext)
		if err := p.App().inject(eda); err != nil {
			p.App().Flash().Err(err)
			return evt
		}
	} else {
		edaEx := NewEnvoyApiRouterView(client.NewGVR("edaEx"))
		edaEx.SetContextFn(p.coContext)
		if err := p.App().inject(edaEx); err != nil {
			p.App().Flash().Err(err)
			return evt
		}
	}
	return nil
}

func (p *Pod) showProxyInfo(evt *tcell.EventKey) *tcell.EventKey {

	path := p.GetTable().GetSelectedItem()
	//log.Info().Msgf("get pods %s in showProxyMetaInfo", path)
	if path == "" {
		return evt
	}

	if has, _ := p.containProxy(); !has {
		log.Debug().Msgf("pod %s has no proxy, skip", path)
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

	if has, _ := p.containProxy(); !has {
		log.Debug().Msgf("pod %s has no proxy, skip", path)
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

func (p *Pod) containProxy() (bool, bool) {
	path := p.GetTable().GetSelectedItem()
	log.Debug().Msgf("get path %s", path)
	hasProxy := false
	sidecar := false

	if path == "" {
		return hasProxy, sidecar
	}
	po, err := p.App().factory.Get("v1/pods", p.GetTable().GetSelectedItem(), true, labels.Everything())
	if err != nil {
		log.Error().Msgf("get pods %s err, %s", p.GetTable().GetSelectedItem(), err.Error())
		return hasProxy, sidecar
	}
	pod := v1.Pod{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(po.(*unstructured.Unstructured).Object, &pod)
	if err != nil {
		log.Error().Msgf("Unstructured convert to pods err, %s", err.Error())
		return hasProxy, sidecar
	}

	name := ""
	for _, co := range pod.Status.ContainerStatuses {
		if co.Name == "istio-proxy" || co.Name == "gateway-proxy" {
			hasProxy = true
			name = co.Name
			break
		}
	}
	if name == "istio-proxy" {
		sidecar = true
	}

	return hasProxy, sidecar
}

func (p *Pod) i9sExtension(evt *tcell.EventKey) *tcell.EventKey {
	sel := p.GetTable().GetSelectedItem()
	log.Debug().Msgf("get sel %s in debug", sel)
	if sel == "" {
		return evt
	}

	has, _ := p.containProxy()
	if !has {
		log.Debug().Msgf("pod %s has no proxy, skip", sel)
		return nil
	}

	env := p.GetTable().envFn()
	log.Debug().Msgf("get env in pod_i9s %+v", env)

	iview := NewI9sExtensionView(client.NewGVR("i9sExtension"))
	iview.SetContextFn(p.coContext)
	if err := p.App().inject(iview); err != nil {
		p.App().Flash().Err(err)
		return evt
	}
	return nil
}
