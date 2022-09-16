package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/derailed/k9s/internal/render"
	"github.com/rs/zerolog/log"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"os/exec"
	"strings"
)

// IstioAdsz
type IstioAdsz struct {
	NonResource
}

// List
func (i IstioAdsz) List(ctx context.Context, ns string) ([]runtime.Object, error) {
	oo := make([]runtime.Object, 0)

	parent, ok := ctx.Value("parent").(string)
	if !ok {
		log.Error().Msgf("Expecting a string but got %T", ctx.Value("parent"))
		return oo, nil
	}
	rev, _ := parse(parent)
	names, namespace := i.getPilotNameNamespace(rev)
	if len(names) == 0 || namespace == "" {
		return oo, nil
	}

	instance2adsz := make(map[string]Adsz, 0)

	for _, name := range names {
		c := getAdsz(namespace,name)
		adsz := parseAdsz(c)
		instance2adsz[name] = adsz
	}

	for name, item := range instance2adsz {
		for _, con := range item.Clients {
			oo = append(oo, render.IstioAdszRes{
				ConnectionId: con.ConnectionId,
				ConnectedAt: con.ConnectedAt,
				Address: con.Address,
				IstioInstance: fmt.Sprintf("%s/%s", namespace, name),
				Parent: strings.Join([]string{namespace, name}, "#") ,
			})
		}
	}
	return oo, nil
}

type Adsz struct {
	TotalClients int32 `json:"totalClients"`
	Clients []Client `json:"clients"`
}

type Client struct {
	ConnectionId string `json:"connectionId"`
	ConnectedAt string `json:"connectedAt"`
	Address string `json:"address"`
	Metadata interface{} `json:"metadata"`
	watches interface{} `json:"watches"`
}

func getAdsz(namespace, name string) string {
	str := fmt.Sprintf("kubectl exec %s -n %s -- curl localhost:15014/debug/adsz -s", name, namespace)
	log.Debug().Msgf("exec %s", str)
	out, err := exec.Command("sh", []string{"-c", str}...).Output()
	if err != nil {
		log.Error().Msgf("exe cmd err, %s", err)
		return ""
	}
	return string(out)
}

func parseAdsz(c string) Adsz {
	adsz := Adsz{}
	if err := json.Unmarshal([]byte(c), &adsz); err != nil {
		log.Error().Msgf("unmarshal adsz to struct %d", err.Error())
		return Adsz{}
	}
	return adsz
}

// 获取全部adsz
func (i *IstioAdsz) getPilotNameNamespace(rev string) ([]string, string) {
	instance := make([]string, 0)
	instanceNs := ""
	label := map[string]string{istioRev: rev, "app": "istiod"}
	oo, err := i.Factory.List("v1/pods", "", true, labels.Set(label).AsSelector())
	if err != nil {
		log.Error().Msgf("list pods with rev:%s, app:istiod err %s", rev, err.Error())
		return instance, instanceNs
	}
	for _, o := range oo {
		var po v1.Pod
		err = runtime.DefaultUnstructuredConverter.FromUnstructured(o.(*unstructured.Unstructured).Object, &po)
		if err != nil {
			return instance, instanceNs
		}
		instance = append(instance, po.Name)
		instanceNs = po.Namespace
	}
	return instance, instanceNs
}

