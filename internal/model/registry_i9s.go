package model

import (
	"github.com/derailed/k9s/internal/dao"
	"github.com/derailed/k9s/internal/render"
)

func init() {
	Registry["istio"] = ResourceMeta{
		DAO:      &dao.Istio{},
		Renderer: &render.Istio{},
	}
	Registry["eda"] = ResourceMeta{
		DAO:      &dao.EnvoyApi{},
		Renderer: &render.EnvoyApi{},
	}
	Registry["edaEx"] = ResourceMeta{
		DAO:      &dao.EnvoyApiRouter{},
		Renderer: &render.EnvoyApiRouter{},
	}
	Registry["ic"] = ResourceMeta{
		DAO:      &dao.IstioConfig{},
		Renderer: &render.IstioConfig{},
	}
	Registry["ida"] = ResourceMeta{
		DAO:      &dao.IstioApi{},
		Renderer: &render.IstioApi{},
	}
	Registry["pilot"] = ResourceMeta{
		DAO:      &dao.IstioPilot{},
		Renderer: &render.IstioPilot{},
	}
	Registry["proxyID"] = ResourceMeta{
		DAO:      &dao.IstioProxyID{},
		Renderer: &render.IstioProxyID{},
	}
	Registry["xps"] = ResourceMeta{
		DAO:      &dao.IstioXdsPushStats{},
		Renderer: &render.IstioXdsPushStats{},
	}
	Registry["configz"] = ResourceMeta{
		DAO:      &dao.IstioConfigz{},
		Renderer: &render.IstioConfigz{},
	}
	Registry["adsz"] = ResourceMeta{
		DAO:      &dao.IstioAdsz{},
		Renderer: &render.IstioAdsz{},
	}
	Registry["proxyInfo"] = ResourceMeta{
		DAO:      &dao.ProxyInfo{},
		Renderer: &render.ProxyInfo{},
	}
	Registry["proxyInfoEx"] = ResourceMeta{
		DAO:      &dao.ProxyInfoEx{},
		Renderer: &render.ProxyInfoEx{},
	}
	Registry["istioctlView"] = ResourceMeta{
		DAO:      &dao.IstioctlView{},
		Renderer: &render.IstioctlView{},
	}
}
