GO_FLAGS   ?=
NAME       := i9s
OUTPUT_BIN ?= execs/${NAME}
PACKAGE    := github.com/derailed/$(NAME)
GIT_REV    ?= $(shell git rev-parse --short HEAD)
BRANCH     ?= $(shell git symbolic-ref --short -q HEAD)
TAG        ?= $(shell git tag --points-at HEAD)
SOURCE_DATE_EPOCH ?= $(shell date +%s)
DATE       ?= $(shell date -u -d @${SOURCE_DATE_EPOCH} +"%Y-%m-%dT%H:%M:%SZ")
REPO       := slimeio
IMG_NAME   := i9s
IMAGE      := ${REPO}/${IMG_NAME}:${TAG}

default: help

test:   ## Run all tests
	@go clean --testcache && go test ./...

cover:  ## Run test coverage suite
	@go test ./... --coverprofile=cov.out
	@go tool cover --html=cov.out

build:  ## Builds the CLI
	@go build ${GO_FLAGS} \
	-ldflags "-w -s -X ${PACKAGE}/cmd.version=${VERSION} -X ${PACKAGE}/cmd.commit=${GIT_REV} -X ${PACKAGE}/cmd.date=${DATE}" \
	-a -tags netgo -o ${OUTPUT_BIN} main.go

kubectl-stable-version:  ## Get kubectl latest stable version
	@curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt

img: build    ## Build Docker Image
	@docker build -t ${IMAGE} .

push: img
	@docker push ${IMAGE}

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":[^:]*?## "}; {printf "\033[38;5;69m%-30s\033[38;5;38m %s\033[0m\n", $$1, $$2}'
