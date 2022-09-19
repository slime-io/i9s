package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/derailed/k9s/internal/render"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/prom2json"
	"github.com/rs/zerolog/log"
	"k8s.io/apimachinery/pkg/runtime"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)
var pre = make(map[string]string)

type IstioXdsPushStats struct {
	NonResource
}

// List
func (i IstioXdsPushStats) List(ctx context.Context, ns string) ([]runtime.Object, error) {
	oo := make([]runtime.Object, 0)
	now := time.Now()
	log.Debug().Msgf("begin time is %s", now.String())

	parent, ok := ctx.Value("parent").(string)
	if !ok {
		log.Error().Msgf("Expecting a string but got %T", ctx.Value("parent"))
		return oo, nil
	}

	m := new(pilot2poryx)
	err := json.Unmarshal([]byte(parent), &m)
	if err != nil {
		log.Error().Msgf("unmarshal from parent %s get err, %s", parent, err)
		return oo, nil
	}
	parts := strings.Split(m.Name, "/")
	if len(parts) != 2 {
		return oo, nil
	}
	str := fetchMetric(parts[1], parts[0])
	metrics := parseMetrics(str)
	log.Debug().Msgf("list metrics %s", metrics)

	res := make(map[string]string)
	if len(pre) == 0 {
		log.Info().Msgf("first enter")
		pre = metrics
	} else {
		for k, v := range metrics {
			newVal := toNum(v)
			oldVal := toNum(pre[k])
			res[k] =  strconv.Itoa(newVal-oldVal)
		}
		pre = metrics
	}

	tmp := []string{"cds", "eds", "lds", "rds"}
	for _, item := range tmp {
		oo = append(oo, render.IstioXdsPushStatsRes{Name: item, Metric: res})
	}
	return oo, nil
}


func parseMetrics(content string) map[string]string {
	mfChan := make(chan *dto.MetricFamily, 1024)
	myReader := strings.NewReader(content)
	err := ParseReader(myReader, mfChan)
	if err != nil {
		fmt.Fprintln(os.Stderr, "read marshaling JSON:", err)
		os.Exit(1)
	}
	m := make(map[string]*prom2json.Family, 0)
	for mf := range mfChan {
		f := prom2json.NewFamily(mf)
		m[f.Name] = f
	}
	metrics := make(map[string]string)
	if push, ok := m["pilot_xds_pushes"]; ok {
		for _, item := range push.Metrics {
			if metric, ok := item.(prom2json.Metric); ok {
				labels := metric.Labels
				if labels["type"] != "" {
					metrics[labels["type"]] = metric.Value
				}
			}
		}
	}
	return metrics
}

func fetchMetric(name, namespace string) string {

	str := fmt.Sprintf("kubectl exec %s -n %s -- curl localhost:15014/metrics -s", name, namespace)
	log.Debug().Msgf("exec %s", str)
	out, err := exec.Command("sh", []string{"-c", str}...).Output()
	if err != nil {
		log.Error().Msgf("execProxyCmd err, %s", err)
		return ""
	}
	return string(out)
}

func toNum(s string) int {

	var num1 float64
	fmt.Sscanf(s, "%e", &num1)
	num2 := fmt.Sprintf("%.f", num1)

	num, err := strconv.Atoi(num2)
	if err != nil {
		return 0
	}
	return num
}