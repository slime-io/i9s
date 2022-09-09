package main

import (
	"flag"
	"runtime"

	"github.com/derailed/k9s/cmd"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/klog/v2"
)

func init() {
	klog.InitFlags(nil)

	if err := flag.Set("logtostderr", "false"); err != nil {
		panic(err)
	}
	if err := flag.Set("alsologtostderr", "false"); err != nil {
		panic(err)
	}
	if err := flag.Set("stderrthreshold", "fatal"); err != nil {
		panic(err)
	}
	if err := flag.Set("one_output", "true"); err != nil {
		panic(err)
	}
	if err := flag.Set("v", "0"); err != nil {
		panic(err)
	}
}

func main() {
	runtime.GOMAXPROCS(2)
	cmd.Execute()
}
