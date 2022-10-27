package dao

import (
	"github.com/derailed/k9s/internal/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func loadi9s(m ResourceMetas) {
	m[client.NewGVR("istio")] = metav1.APIResource{
		Name:         "istio",
		Kind:         "Istio",
		SingularName: "istio",
		Verbs:        []string{},
		Categories:   []string{"k9s"},
	}
	m[client.NewGVR("eda")] = metav1.APIResource{
		Name:         "eda",
		Kind:         "eda",
		SingularName: "eda",
		Verbs:        []string{},
		Categories:   []string{"k9s"},
	}
	m[client.NewGVR("ida")] = metav1.APIResource{
		Name:         "ida",
		Kind:         "ida",
		SingularName: "ida",
		Verbs:        []string{},
		Categories:   []string{"k9s"},
	}
	m[client.NewGVR("ic")] = metav1.APIResource{
		Name:         "ic",
		Kind:         "ic",
		SingularName: "ic",
		Verbs:        []string{},
		Categories:   []string{"k9s"},
	}
	m[client.NewGVR("pilot")] = metav1.APIResource{
		Name:         "pilot",
		Kind:         "pilot",
		SingularName: "pilot",
		Verbs:        []string{},
		Categories:   []string{"k9s"},
	}
	m[client.NewGVR("proxyID")] = metav1.APIResource{
		Name:         "proxyID",
		Kind:         "proxyID",
		SingularName: "proxyID",
		Verbs:        []string{},
		Categories:   []string{"k9s"},
	}
	m[client.NewGVR("xps")] = metav1.APIResource{
		Name:         "xps",
		Kind:         "xps",
		SingularName: "xps",
		Verbs:        []string{},
		Categories:   []string{"k9s"},
	}
	m[client.NewGVR("configz")] = metav1.APIResource{
		Name:         "configz",
		Kind:         "configz",
		SingularName: "configz",
		Verbs:        []string{},
		Categories:   []string{"k9s"},
	}
	m[client.NewGVR("adsz")] = metav1.APIResource{
		Name:         "adsz",
		Kind:         "adsz",
		SingularName: "adsz",
		Verbs:        []string{},
		Categories:   []string{"k9s"},
	}
	m[client.NewGVR("proxyInfo")] = metav1.APIResource{
		Name:         "proxyInfo",
		Kind:         "proxyInfo",
		SingularName: "proxyInfo",
		Verbs:        []string{},
		Categories:   []string{"k9s"},
	}
	m[client.NewGVR("proxyInfoEx")] = metav1.APIResource{
		Name:         "proxyInfoEx",
		Kind:         "proxyInfoEx",
		SingularName: "proxyInfoEx",
		Verbs:        []string{},
		Categories:   []string{"k9s"},
	}
	m[client.NewGVR("istioctlView")] = metav1.APIResource{
		Name:         "istioctlView",
		Kind:         "istioctlView",
		SingularName: "istioctlView",
		Verbs:        []string{},
		Categories:   []string{"k9s"},
	}
}
