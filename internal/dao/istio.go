package dao

import (
	"context"
	"github.com/derailed/k9s/internal/render"
	"github.com/rs/zerolog/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	istioRev = "istio.io/rev"
)

type Istio struct {
	NonResource
}

func (i *Istio) List(ctx context.Context, ns string) ([]runtime.Object, error) {
	oo := make([]runtime.Object, 0)
	revs := i.GetRev()
	for _, f := range revs {
		oo = append(oo, render.IstioRes{Name: f})
	}
	return oo, nil
}

// Describe returns the revision notes.
func (i *Istio) Describe(path string) (string, error) {
	return "", nil
}

// ToYAML returns the revision manifest.
func (i *Istio) ToYAML(path string, showManaged bool) (string, error) {
	return "", nil
}

func (i *Istio) GetRev() []string {
	var revisions []string
	revs := make(map[string]string, 0)

	dial, err := i.Client().Dial()
	if err != nil {
		log.Error().Msgf("get client err in dao/istio, %s", err)
		return revisions
	}
	dps, err := dial.AppsV1().Deployments("").
		List(context.TODO(), metav1.ListOptions{
			LabelSelector: toSelector(map[string]string{"app": "istiod"}),
		})

	if err != nil {
		log.Error().Msgf("list dps err %s", err.Error())
		return revisions
	}

	for _, dp := range dps.Items {
		if rev, ok := dp.Labels[istioRev]; ok {
			if _, exist := revs[rev]; !exist {
				revisions = append(revisions, rev)
			}
		}
	}
	return revisions
}
