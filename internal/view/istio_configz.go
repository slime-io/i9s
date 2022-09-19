package view

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/derailed/k9s/internal/client"
	"github.com/derailed/k9s/internal/render"
	"github.com/derailed/k9s/internal/ui"
	"github.com/gdamore/tcell/v2"
	"github.com/rs/zerolog/log"
	"github.com/yudai/gojsondiff"
	"github.com/yudai/gojsondiff/formatter"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"os/exec"
	"strings"
	"time"
)

type IstioConfigzView struct {
	ResourceViewer
}

func NewIstioConfigzView(gvr client.GVR) ResourceViewer {
	c := IstioConfigzView{
		ResourceViewer: NewBrowser(gvr),
	}
	c.GetTable().SetColorerFn(render.IstioConfig{}.ColorerFunc())
	c.GetTable().SetBorderFocusColor(tcell.ColorMediumSpringGreen)
	c.GetTable().SetSelectedStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorMediumSpringGreen).Attributes(tcell.AttrNone))
	c.SetContextFn(c.chartContext)
	c.GetTable().SetEnterFn(c.enter)
	c.AddBindKeysFn(c.bindKeys)
	return &c
}

func (i *IstioConfigzView) chartContext(ctx context.Context) context.Context {
	return ctx
}

func (i *IstioConfigzView) enter(app *App, model ui.Tabular, gvr, path string) {

	kind, namespace, name, podnamespace, podname, err := parseNNK(path)
	if err != nil {
		log.Error().Msgf("get err in  %s", err)
		return
	}
	b := content(kind, namespace, name, podnamespace, podname)

	details := NewDetails(i.App(), "json", "info", true)
	details.Update(string(b))
	if err := i.App().inject(details); err != nil {
		i.App().Flash().Err(err)
	}
}

func (i *IstioConfigzView) bindKeys(aa ui.KeyActions) {
	aa.Add(ui.KeyActions{ui.KeyW: ui.NewKeyAction("watch", i.watch, true)})
	aa.Add(ui.KeyActions{ui.KeyD: ui.NewKeyAction("diff", i.diff, true)})
}

func (i * IstioConfigzView) diff(evt *tcell.EventKey) *tcell.EventKey {

	path := i.GetTable().GetSelectedItem()
	if path == "" {
		return evt
	}
	// get config in configz
	kind, namespace, name, podnamespace, podname, err := parseNNK(path)
	if err != nil {
		log.Error().Msgf("get err in IstioConfigzView parseNNK %s", err)
		return evt
	}
	configInIstio := content(kind, namespace, name, podnamespace, podname)

	kk := strings.ToLower(kind)
	if strings.HasSuffix(kk, "y") {
		kk = kk[:len(kk)-1]
		kk = kk + "ies"
	} else {
		kk = kk + "s"
	}

	// get config in k8s
	gvr := fmt.Sprintf("networking.istio.io/v1alpha3/%s", kk)
	nn := fmt.Sprintf("%s/%s", strings.ToLower(namespace), strings.ToLower(name))

	o, err := i.App().factory.Get(gvr, nn, true, labels.NewSelector())
	if err != nil {
		log.Error().Msgf("get gvr %s, nn %s err: %s", gvr, nn, err.Error())
		return evt
	}

	// convert k8s config to configZ
	c := Configz{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(o.(*unstructured.Unstructured).Object, &c)
	if err != nil {
		log.Error().Msgf("convert runtime object to istio config err", err.Error())
		return evt
	}
	configInKubernetst, err:= json.MarshalIndent(c, "", "  ")
	if err != nil {
		log.Error().Msgf("marshal runtime object to json err", err)
		return evt
	}

	// diff
	differ := gojsondiff.New()
	diff, err:= differ.Compare(configInKubernetst, configInIstio)
	if err != nil {
		log.Error().Msgf("get err in json diff, %s", err)
		return evt
	}

	var aJson map[string]interface{}
	json.Unmarshal(configInKubernetst, &aJson)

	config := formatter.AsciiFormatterConfig{
		ShowArrayIndex: true,
		Coloring:       false,
	}
	formatter := formatter.NewAsciiFormatter(aJson, config)
	res, err := formatter.Format(diff)

	//format := formatter.NewDeltaFormatter()
	//res, err := format.Format(diff)
	if err != nil {
		log.Error().Msgf("get err in json diff, %s", err)
		return evt
	}

	details := NewDetails(i.App(), "diff", "k8s->istio", true)
	details.Update(res)
	if err := i.App().inject(details); err != nil {
		i.App().Flash().Err(err)
	}
	return evt
}

func (i * IstioConfigzView) watch(evt *tcell.EventKey) *tcell.EventKey {
	path := i.GetTable().GetSelectedItem()
	log.Info().Msgf("get path %s in IstioConfigzView watch", path)
	if path == "" {
		return evt
	}

	kind, namespace, name, podnamespace, podname, err := parseNNK(path)
	if err != nil {
		log.Error().Msgf("get err in IstioConfigzView parseNNK %s", err)
		return nil
	}

	old := content(kind, namespace, name, podnamespace, podname)
	details := NewDetails(i.App(), "json", "info", true)
	details.Update(string(old))
	if err := i.App().inject(details); err != nil {
		i.App().Flash().Err(err)
	}
	go func() {
		timer := time.NewTicker(5*time.Second)
		buf := make([]byte,0)
		prefix := ""
		for {
			select {
			case <- timer.C:
				if details != nil {
					watched := content(kind, namespace, name, podnamespace, podname)
					b, changed, err := jsonDiff(old, watched)
					if err != nil {
						log.Error().Msgf("get err in json diff, %s", err.Error())
						continue
					}
					if changed {
						prefix = fmt.Sprintf("## config has changed in last flush, time: %s\n", time.Now().Format(time.RFC3339))
						buf = []byte(prefix)
						buf = append(buf, b...)
						details.Update(string(buf))
						prefix = ""
						old = watched
					}
				}
			}
		}
	}()
	return nil
}



func content(kind, namespace, name, podnamespace, podname string ) []byte {

	var b []byte
	conifgz := getConfigz(podnamespace, podname)
	if conifgz == "" {
		return b
	}
	items := parseConfigz(conifgz)
	if len(items) == 0 {
		return b
	}
	item := parse(kind, namespace, name, items)
	b, err := json.MarshalIndent(item, "", "  ")
	if err != nil {
		log.Error().Msgf("marshal configz err %s", err.Error())
		return []byte{}
	}
	return b
}


func parseNNK(path string) (string, string, string, string, string, error) {
	parts := strings.Split(path, "#")
	if len(parts) != 5 {
		return "", "", "", "", "", fmt.Errorf("invalid path %s", path)
	}
	return parts[0], parts[1], parts[2], parts[3], parts[4], nil
}

type Configz struct {
	Kind string `json:"kind"`
	Metadata MetaData `json:"metadata"`
	ApiVersion string `json:"apiVersion"`
	Spec interface{} `json:"spec"`
}

type MetaData struct {
	Name string `json:"name"`
	Namespace string `json:"namespace"`
	Labels interface{} `json:"labels"`
	Annotations interface{} `json:"annotations"`
	ResourceVersion string `json:"resourceVersion"`
	CreationTimestamp string `json:"creationTimestamp,"`
	Generation int64 `json:"generation,omitempty"`
}

func getConfigz(namespace, name string) string {
	str := fmt.Sprintf("kubectl exec %s -n %s -- curl localhost:15014/debug/configz -s", name, namespace)
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


func parse(kind, namespace, name string, items []Configz) Configz {
	var config Configz
	for _, item := range items {
		if item.Kind != kind {
			continue
		}
		if item.Metadata.Namespace != namespace {
			continue
		}
		if item.Metadata.Name != name {
			continue
		}
		config = item
		break
	}
	return config
}


func jsonDiff(s1, s2 []byte) (string, bool, error){
	differ := gojsondiff.New()
	diff, err:= differ.Compare(s1, s2)
	if err != nil {
		return "", false , err
	}

	var aJson map[string]interface{}
	json.Unmarshal(s1, &aJson)

	config := formatter.AsciiFormatterConfig{
		ShowArrayIndex: true,
		Coloring:       false,
	}
	formatter := formatter.NewAsciiFormatter(aJson, config)
	b, err := formatter.Format(diff)
	return b, diff.Modified(), err
}

// istio config
