package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/derailed/k9s/internal"
	"github.com/derailed/k9s/internal/render"
	"github.com/rs/zerolog/log"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"os/exec"
	"strings"
)

// IstioConfig IC
type IstioConfigz struct {
	NonResource
}

// List
func (i IstioConfigz) List(ctx context.Context, ns string) ([]runtime.Object, error) {
	oo := make([]runtime.Object, 0)

	parent, ok := ctx.Value(internal.Parent).(string)
	if !ok {
		log.Error().Msgf("Expecting a string but got %T", ctx.Value(internal.Parent))
		return oo, nil
	}
	rev, _ := parse(parent)
	namespace, name := i.getPilotNameNamespace(rev)
	if name == "" || namespace == "" {
		return oo, nil
	}
	c := getConfigz(namespace, name)
	items := parseConfigz(c)
	if len(items) > 0 {
		for _, item := range items {
			oo = append(oo, render.IstioConfigzRes{
				Name:      item.Metadata.Name,
				Namespace: item.Metadata.Namespace,
				Kind:      item.Kind,
				Parent:    strings.Join([]string{namespace, name}, "#"),
			})
		}
	}
	return oo, nil
}

type Configz struct {
	Kind     string   `json:"kind"`
	Metadata MetaData `json:"metadata"`
}

type MetaData struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

func getConfigz(namespace, name string) string {
	str := fmt.Sprintf("kubectl exec %s -n %s -- curl localhost:15014/debug/configz -s", name, namespace)
	log.Debug().Msgf("exec %s", str)
	out, err := exec.Command("sh", []string{"-c", str}...).Output()
	if err != nil {
		log.Error().Msgf("execProxyCmd err, %s", err)
		return ""
	}
	return string(out)
}

func parseConfigz(c string) []Configz {
	configz := make([]Configz, 0)
	if err := json.Unmarshal([]byte(c), &configz); err != nil {
		log.Error().Msgf("unmarshal configz to struct %d", err.Error())
		return []Configz{}
	}
	return configz
}

func (i *IstioConfigz) getPilotNameNamespace(rev string) (string, string) {
	label := map[string]string{istioRev: rev, "app": "istiod"}
	oo, err := i.Factory.List("v1/pods", "", true, labels.Set(label).AsSelector())
	if err != nil {
		log.Error().Msgf("list pods with rev:%s, app:istiod err %s", rev, err.Error())
		return "", ""
	}
	for _, o := range oo {
		var po v1.Pod
		err = runtime.DefaultUnstructuredConverter.FromUnstructured(o.(*unstructured.Unstructured).Object, &po)
		if err != nil {
			return "", ""
		}
		return po.Namespace, po.Name
	}
	return "", ""
}
