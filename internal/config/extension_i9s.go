package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"strings"
)

var ViewToExtension = make(map[string][]Extension)

// usage:

//var LocalExtensions = `extensions:
//- view: istioView
//  name: rev-reset
//  command: sh
//  args:
//  - "cmd reset--rev $ISTIO_REV"
//- view: istioView
//  name: echo
//  command: sh
//  args:
//  - "echo $ISTIO_REV"
//- view: sidecarView
//  name: describe
//  command: sh
//  args:
//  - "kubectl describe pods $NAME -n $NAMESPACE"
//`

var LocalExtensions = ``

type Extensions struct {
	Extension []Extension `yaml:"extensions"`
}

type Extension struct {
	View    string   `yaml:"view"`
	Name    string   `yaml:"name"`
	Args    []string `yaml:"args"`
	Command string   `yaml:"command"`
}

func (p Extension) String() string {
	return fmt.Sprintf("[%s-%s] %s(%s)", p.View, p.Name, p.Command, strings.Join(p.Args, " "))
}

func init() {
	LoadExtensions()
}

func LoadExtensions() {
	var pp Extensions
	if err := yaml.Unmarshal([]byte(LocalExtensions), &pp); err != nil {
		return
	}
	for _, v := range pp.Extension {
		if _, ok := ViewToExtension[v.View]; !ok {
			ViewToExtension[v.View] = make([]Extension, 0)
		}
		ViewToExtension[v.View] = append(ViewToExtension[v.View], v)
	}
}
