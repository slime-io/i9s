module github.com/derailed/k9s

go 1.17

replace (
	github.com/docker/distribution => github.com/docker/distribution v0.0.0-20191216044856-a8371794149d
	github.com/docker/docker => github.com/moby/moby v17.12.0-ce-rc1.0.20200618181300-9dc6525e6118+incompatible
	github.com/gdamore/tcell/v2 => github.com/derailed/tcell/v2 v2.3.1-rc.2
)

require (
	github.com/adrg/xdg v0.4.0
	github.com/atotto/clipboard v0.1.4
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/cenkalti/backoff/v4 v4.1.2
	github.com/derailed/popeye v0.9.8
	github.com/derailed/tview v0.6.6
	github.com/fatih/color v1.13.0
	github.com/fsnotify/fsnotify v1.5.1
	github.com/fvbommel/sortorder v1.0.2
	github.com/gdamore/tcell/v2 v2.4.0
	github.com/ghodss/yaml v1.0.0
	github.com/mattn/go-runewidth v0.0.13
	github.com/petergtz/pegomock v2.9.0+incompatible
	github.com/rakyll/hey v0.1.4
	github.com/rs/zerolog v1.26.0
	github.com/sahilm/fuzzy v0.1.0
	github.com/spf13/cobra v1.2.1
	github.com/stretchr/testify v1.7.0
	golang.org/x/text v0.3.7
	gopkg.in/yaml.v2 v2.4.0
	helm.sh/helm/v3 v3.7.1
	k8s.io/api v0.22.3
	k8s.io/apimachinery v0.22.3
	k8s.io/cli-runtime v0.22.3
	k8s.io/client-go v0.22.3
	k8s.io/klog/v2 v2.30.0
	k8s.io/kubectl v0.22.3
	k8s.io/metrics v0.22.3
	sigs.k8s.io/yaml v1.3.0
)

require (
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/mattn/go-colorable v0.1.9
	github.com/prometheus/client_model v0.2.0
	github.com/prometheus/prom2json v1.3.1
)
