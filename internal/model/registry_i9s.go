package model

import (
	"github.com/derailed/k9s/internal/dao"
	"github.com/derailed/k9s/internal/render"
)

func init() {
	Registry["istio"] = ResourceMeta{
		DAO:          &dao.Istio{},
		Renderer:     &render.Istio{},
	}
	Registry["eda"] = ResourceMeta{
		DAO:          &dao.EnvoyApi{},
		Renderer:     &render.EnvoyApi{},
	}
	Registry["ic"] = ResourceMeta{
		DAO:          &dao.IstioConfig{},
		Renderer:     &render.IstioConfig{},
	}
	Registry["ida"] = ResourceMeta{
		DAO:          &dao.IstioApi{},
		Renderer:     &render.IstioApi{},
	}
	Registry["pilot"] = ResourceMeta{
		DAO:          &dao.IstioPilot{},
		Renderer:     &render.IstioPilot{},
	}
	Registry["proxyID"] = ResourceMeta{
		DAO:          &dao.IstioProxyID{},
		Renderer:     &render.IstioProxyID{},
	}
	Registry["xps"] = ResourceMeta{
		DAO:          &dao.IstioXdsPushStats{},
		Renderer:     &render.IstioXdsPushStats{},
	}
	Registry["configz"] = ResourceMeta{
		DAO:          &dao.IstioConfigz{},
		Renderer:     &render.IstioConfigz{},
	}
	Registry["adsz"] = ResourceMeta{
		DAO:          &dao.IstioAdsz{},
		Renderer:     &render.IstioAdsz{},
	}
	Registry["iptableInfo"] = ResourceMeta{
		DAO:          &dao.IptableInfo{},
		Renderer:     &render.IptableInfo{},
	}
}