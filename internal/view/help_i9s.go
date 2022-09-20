package view

import (
	"encoding/json"
	"fmt"
	"github.com/derailed/k9s/internal/dao"
	"github.com/derailed/k9s/internal/port"
	"github.com/derailed/k9s/internal/watch"
	"github.com/rs/zerolog/log"
	"io"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/portforward"
	"net"
	"net/http"
	"sigs.k8s.io/yaml"
	"strconv"
	"strings"
	"time"
)
// i9s help

const (
	istioRev = "istio.io/rev"
)

var (
	apisNeedPodSelect = []string{"adsz", "connections", "instancesz", "syncz", "edsz", "sidecarz", "config_dump", "metrics", "xds_push_stats"}
	apisNeedProxyID = []string{"edsz", "sidecarz", "config_dump"}
	exApis = []string{"configzEx", "adszEx"}
)

func execi9sCmd(i ResourceViewer, path, podName, rev, ProxyID string){
	// if podName is "", choose first one in rev
	if podName == "" {
		podName = getPod(i, rev, "")
	}

	// check the port-forward
	if _, ok := i.App().factory.ForwarderFor(fwFQN(podName, container)); ok {
		log.Info().Msgf("A port-forward already exist on pod|container %s|discovery", podName)
		return
	}

	allocatePort, err := availablePort("localhost")
	if err != nil {
		log.Info().Msgf("failure allocating port: %v", err)
		return
	}

	nodePort := strconv.Itoa(allocatePort)
	tt, err := port.ToTunnels("localhost", fmt.Sprintf("%s::%s", container, containerPort), nodePort)
	if err != nil {
		log.Error().Msgf("port convert to ToTunnels get err %s", err)
		return
	}

	// run port-forward
	if err = RunI9sPortForward(i, tt, podName); err != nil {
		log.Error().Msgf("get error in RunPortForward, %s", err)
		return
	}

	id := dao.PortForwardID(fmt.Sprintf("%s|%s", podName, container), container, fmt.Sprintf("%s:%s", nodePort, containerPort))
	defer i.App().factory.DeleteForwarder(id)

	<- time.After(200*time.Millisecond)

	// exec cmd
	con, err := net.Dial("tcp", fmt.Sprintf("localhost:%s", nodePort))
	if err == nil {
		con.Close()
	} else {
		<- time.After(400*time.Millisecond)
	}
	execCmd(i, path, nodePort, ProxyID)
}

func execI9sCmdWithHttp(i ResourceViewer, path, podName, rev string) []string {
	var ids []string
	if podName == "" {
		podName = getPod(i, rev, "")
	}
	if _, ok := i.App().factory.ForwarderFor(fwFQN(podName, container)); ok {
		log.Info().Msgf("A port-forward already exist on pod|container %s|discovery", podName)
		return ids
	}
	allocatePort, err := availablePort("localhost")
	if err != nil {
		log.Info().Msgf("failure allocating port: %v", err)
		return ids
	}
	localPort := strconv.Itoa(allocatePort)
	tt, err := port.ToTunnels("localhost", fmt.Sprintf("%s::%s", container, containerPort), localPort)
	if err != nil {
		log.Error().Msgf("port convert to ToTunnels get err %s", err)
		return ids
	}
	if err = RunI9sPortForward(i, tt, podName); err != nil {
		log.Error().Msgf("get error in RunPortForward, %s", err)
		return ids
	}
	defer i.App().factory.DeleteForwarder(dao.PortForwardID(fmt.Sprintf("%s|%s", podName, container), container, fmt.Sprintf("%s:%s", localPort, containerPort)))
	<- time.After(100*time.Millisecond)
	ids = execHttp(path, localPort)
	return ids
}

func execHttp(path, localPort string)  []string {
	var ids []string
	url := fmt.Sprintf("http://localhost:%s/debug/connections", localPort)
	log.Info().Msgf("request url with http client %s", url)
	res, err := http.Get(url)
	if err != nil {
		log.Error().Msgf("get err when request %s, %s", url, err)
		return ids
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error().Msgf("get err in readall %s", err)
		return ids
	}
	log.Debug().Msgf("get info from pilot connection %s", string(b))
	clients :=  &Clients{}
	err = json.Unmarshal(b, clients)
	if err != nil {
		log.Error().Msgf("unmarshal err in do http to connections, %s", err)
		return ids
	}
	for _, item := range clients.Connected {
		ids = append(ids, item.ConnectionID)
	}
	return ids
}

func RunI9sPortForward(i ResourceViewer, pts port.PortTunnels, podName string) error {
	if err := pts.CheckAvailable(); err != nil {
		return err
	}

	tt := make([]string, 0, len(pts))
	for _, pt := range pts {
		path := fmt.Sprintf("%s|%s", podName, container)
		if _, ok := i.App().factory.ForwarderFor(dao.PortForwardID(path, pt.Container, pt.PortMap())); ok {
			return fmt.Errorf("A port-forward is already active on pod %s", podName)
		}
		pf := dao.NewPortForwarder(i.App().factory)
		fwd, err := pf.Start(path, pt)
		if err != nil {
			return err
		}
		log.Info().Msgf(">>> Starting port forward %q -- %#v", pf.ID(), pt)
		go portForward(i, pf, fwd)
		tt = append(tt, pt.ContainerPort)
	}

	if len(tt) == 1 {
		log.Info().Msgf("PortForward activated %s", tt[0])
		return nil
	}
	return nil
}

func portForward(i ResourceViewer, pf watch.Forwarder, f *portforward.PortForwarder) {
	i.App().factory.AddForwarder(pf)
	pf.SetActive(true)
	if err := f.ForwardPorts(); err != nil {
		log.Info().Msgf("runPortForward err %s", err)
		return
	}
}

func getPod(i ResourceViewer,rev, ns string) string {
	podNames := getAllPod(i, rev, ns)
	log.Info().Msgf("get pilot %v pod in rev %s", podNames, rev)
	if len(podNames) > 0 {
		return podNames[0]
	}
	return ""
}

func getAllPod(i ResourceViewer, rev, ns string) []string {

	podNames := make([]string, 0)
	label := map[string]string{istioRev: rev, "app": "istiod"}
	oo, err := i.App().factory.List("v1/pods", ns,false, labels.Set(label).AsSelector())
	if err != nil {
		return podNames
	}
	for _, o := range oo {
		var po v1.Pod
		err = runtime.DefaultUnstructuredConverter.FromUnstructured(o.(*unstructured.Unstructured).Object, &po)
		if err != nil {
			return podNames
		}
		if err := checkRunningStatus(container, po.Status.ContainerStatuses); err == nil {
			podNames = append(podNames, fqn(po.Namespace, po.Name))
		}
	}
	return podNames
}

func execCmd(i ResourceViewer, path, localPort, ProxyID string){
	cmd := buildCmd(path, localPort, ProxyID)
	cb := func() {
		opts := shellOpts{
			clear:      false,
			binary:     "sh",
			background: false,
			args:       []string{"-c", cmd},
		}
		if run(i.App(), opts) {
			i.App().Flash().Info("command successfully!")
			return
		}
		log.Error().Msgf("command %s failed!", cmd)
	}
	cb()
}

func buildCmd(path, nodePort, proxyID string) string {

	if path == "metrics" {
		cmd:= fmt.Sprintf("curl localhost:%s/%s -s | less", nodePort, path)
		return cmd
	}

	url := fmt.Sprintf("curl localhost:%s/debug/%s", nodePort, path)
	if proxyID != "" {
		proxy := fmt.Sprintf("?proxyID=%s", proxyID)
		url = url + proxy
	}
	cmd := url + " -s | jq . | less"
	log.Info().Msgf("build cmd %s", cmd)
	return cmd
}

func availablePort(localAddr string) (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", net.JoinHostPort(localAddr, "0"))
	if err != nil {
		return 0, err
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	port := l.Addr().(*net.TCPAddr).Port
	return port, l.Close()
}

func needPodSelected(path string) bool {
	for _, item := range apisNeedPodSelect {
		if item == path {
			return true
		}
	}
	return false
}

func needProxyID(path string) bool {
	for _, item := range apisNeedProxyID {
		if item == path {
			return true
		}
	}
	return false
}

func needEx(path string) bool {
	for _, item := range exApis {
		if item == path {
			return true
		}
	}
	return false
}

// istio/config_dump
func formatIstioAPI(api string) (string, error) {
	parts := strings.Split(api, "/")
	if len(parts) != 2 {
		return "" , fmt.Errorf("except 2 items in %s", api)
	}
	return parts[1], nil
}

// 112#istio/config_dump
func parseRevWithAPI(s string) (string, string, error) {
	parts := strings.Split(s, "#")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("except 2 items in %s", s)
	}
	return parts[0], parts[1], nil
}

// 112#configmap
func parseIstioConfig(s string) (string, string, error) {
	parts := strings.Split(s, "#")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("except 2 items in %s", s)
	}
	return parts[0], parts[1], nil
}

//112 # istio/sidecarz # istio-system/istiod-112-79dd58f89-slrgd
func parseIstioPodView(s string) (string, string, string, error) {
	parts := strings.Split(s, "#")
	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("except 3 items in %s", s)
	}
	return parts[0], parts[1], parts[2], nil
}

// 112 # istio/sidecarz # istio-system/istiod-112-79dd58f89-slrgd # a-v1.s-568669c495-4pdvp.powerful-14168
func parseIstioProxyIDViewID(s string) (string, string, string){
	parts := strings.Split(s, "#")
	_, api, pilot, proxy := parts[0], parts[1], parts[2], parts[3]
	return api, pilot, proxy
}

func convert2String(obj interface{}) (string, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		log.Error().Msgf("marshal istio deployment manifest to json err, %s", err)
		return "", err
	}
	data, err = yaml.JSONToYAML(data)
	if err != nil {
		log.Error().Msgf("marshal istio deployment manifest to yaml err, %s", err)
		return "", err
	}
	return string(data), nil
}