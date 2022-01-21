MAKEFLAGS += --warn-undefined-variables
SHELL := /bin/bash -o pipefail
DOCKER_REPO ?= "a3s"
DOCKER_IMAGE ?= "a3s"
DOCKER_TAG ?= "dev"

export GO111MODULE = on

default: codegen sec test a3s cli
.PHONY: ui

lint:
	golangci-lint run \
		--timeout 5m \
		--disable-all \
		--exclude-use-default=false \
		--enable=errcheck \
		--enable=goimports \
		--enable=ineffassign \
		--enable=revive \
		--enable=unused \
		--enable=structcheck \
		--enable=staticcheck \
		--enable=varcheck \
		--enable=deadcode \
		--enable=unconvert \
		--enable=misspell \
		--enable=prealloc \
		--enable=nakedret \
		--enable=typecheck \
		--enable=nilerr \
		./...

test:
	go test ./... -race -cover -covermode=atomic -coverprofile=unit_coverage.cov
	gocov convert ./unit_coverage.cov | gocov-xml > ./coverage.xml

sec:
	gosec -quiet ./...

generate:
	go generate ./...

api:
	cd pkgs/api && make codegen

ui:
	cd ui/login && yarn && yarn build

codegen: api ui generate

a3s:
	cd cmd/a3s && CGO_ENABLED=0 go build -ldflags="-w -s" -trimpath

a3s_linux:
	cd cmd/a3s && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -trimpath

cli:
	cd cmd/a3sctl && CGO_ENABLED=0 go install -ldflags="-w -s" -trimpath

cli_linux:
	cd cmd/a3sctl && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install -ldflags="-w -s" -trimpath

docker: codegen generate a3s_linux package_ca_certs
	mkdir -p docker/in
	cp cmd/a3s/a3s docker/in
	cd docker && docker build -t ${DOCKER_REPO}/${DOCKER_IMAGE}:${DOCKER_TAG} .

package_ca_certs:
	mkdir -p docker/in
	go install github.com/agl/extract-nss-root-certs@latest
	curl -s https://hg.mozilla.org/mozilla-central/raw-file/tip/security/nss/lib/ckfw/builtins/certdata.txt -o certdata.txt
	mkdir -p docker/in
	extract-nss-root-certs > docker/in/ca-certificates.pem
	rm -f certdata.txt
