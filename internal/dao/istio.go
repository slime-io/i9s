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

// Istio represents a istio command.
type Istio struct {
	NonResource
}

func (i *Istio) List(ctx context.Context, ns string) ([]runtime.Object, error) {
	oo := make([]runtime.Object, 0)
	revisions := i.GetRev()
	for _, f := range revisions {
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
	dial, err := i.Client().Dial()
	if err != nil {
		log.Error().Msgf("get client err in dao/istio, %s", err)
		return revisions
	}
	hooks, err := dial.AdmissionregistrationV1beta1().MutatingWebhookConfigurations().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Error().Msgf("list mutatingwebhookconfigurations err %s", err.Error())
		return revisions
	}
	for _, hook := range hooks.Items {
		if rev, ok := hook.GetLabels()[istioRev]; ok {
			revisions = append(revisions, rev)
		}
	}
	log.Info().Msgf("get all istio revs %s", revisions)
	return revisions
}
