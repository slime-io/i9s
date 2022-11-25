package dao

import (
	"context"
	"github.com/derailed/k9s/internal"
	"github.com/derailed/k9s/internal/render"
	"github.com/rs/zerolog/log"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"strings"
)

// IstioPilot IC
type IstioPilot struct {
	NonResource
}

// List
func (i IstioPilot) List(ctx context.Context, ns string) ([]runtime.Object, error) {
	oo := make([]runtime.Object, 0)
	//   112#istio/sidecarz
	parent, ok := ctx.Value(internal.Parent).(string)
	if !ok {
		log.Error().Msgf("Expecting a string but got %T", ctx.Value(internal.Parent))
		return oo, nil
	}
	// log.Info().Msgf("get parent %s in dao istio_pilot", parent)
	rev, api := parse(parent)
	instances, err := i.pilotInstance(rev)
	if err != nil {
		log.Error().Msgf("fetch pilot instance err, %s", err)
		return oo, nil
	}
	for _, item := range instances {
		oo = append(oo, render.IstioPilotRes{Name: item, Revision: rev, Api: api, Parent: parent})
	}
	return oo, nil
}

func (i *IstioPilot) pilotInstance(rev string) ([]string, error) {
	var ins []string
	label := map[string]string{istioRev: rev, "app": "istiod"}
	oo, err := i.Factory.List("v1/pods", "", true, labels.Set(label).AsSelector())
	if err != nil {
		return nil, err
	}
	for _, o := range oo {
		var po v1.Pod
		err = runtime.DefaultUnstructuredConverter.FromUnstructured(o.(*unstructured.Unstructured).Object, &po)
		if err != nil {
			return nil, err
		}
		ins = append(ins, FQN(po.Namespace, po.Name))
	}
	return ins, nil
}

func parse(s string) (string, string) {
	parts := strings.Split(s, "#")
	return parts[0], parts[1]
}
