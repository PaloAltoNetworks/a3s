MAKEFLAGS += --warn-undefined-variables
SHELL := /bin/bash -o pipefail
CONTAINER_ENGINE ?= docker
CONTAINER_REPO ?= "a3s"
CONTAINER_IMAGE ?= "a3s"
CONTAINER_TAG ?= "dev"

export GO111MODULE = on

default: lint vuln test a3s cli
.PHONY: ui docker

## Tests

lint:
	golangci-lint run \
		--timeout=5m \
		--disable-all \
		--exclude-use-default=false \
		--exclude=dot-imports \
		--exclude=package-comments \
		--exclude=unused-parameter \
		--exclude=dot-imports \
		--enable=errcheck \
		--enable=goimports \
		--enable=ineffassign \
		--enable=revive \
		--enable=unused \
		--enable=staticcheck \
		--enable=unconvert \
		--enable=misspell \
		--enable=prealloc \
		--enable=nakedret \
		--enable=typecheck \
		--enable=unparam \
		--enable=gosimple \
		--enable=nilerr \
		./...


test:
	go test ./... -race -cover -covermode=atomic -coverprofile=unit_coverage.out

sec:
	gosec -quiet ./...

vuln:
	govulncheck ./...


## Code generation

generate:
	go generate ./...

api:
	cd pkgs/api && make codegen

ui:
	cd internal/ui/js/login && yarn && yarn build

codegen: api ui generate


## Main build

a3s:
	cd cmd/a3s && CGO_ENABLED=0 go build -ldflags="-w -s" -trimpath

a3s_linux:
	cd cmd/a3s && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -trimpath

cli:
	cd cmd/a3sctl && CGO_ENABLED=0 go install -ldflags="-w -s" -trimpath

cli_linux:
	cd cmd/a3sctl && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install -ldflags="-w -s" -trimpath


## Containers

docker:
	CONTAINER_ENGINE=docker make container

podman:
	CONTAINER_ENGINE=podman make container

container: codegen generate a3s_linux package_ca_certs
	mkdir -p docker/in
	cp cmd/a3s/a3s docker/in
	cd docker && ${CONTAINER_ENGINE} build -t ${CONTAINER_REPO}/${CONTAINER_IMAGE}:${CONTAINER_TAG} .

package_ca_certs:
	mkdir -p docker/in
	go install github.com/agl/extract-nss-root-certs@latest
	curl -s https://hg.mozilla.org/mozilla-central/raw-file/tip/security/nss/lib/ckfw/builtins/certdata.txt -o certdata.txt
	mkdir -p docker/in
	extract-nss-root-certs > docker/in/ca-certificates.pem
	rm -f certdata.txt

# tag the commit, set GITHUB_TOKEN, then run...
release:
	unset GITLAB_TOKEN && goreleaser check && goreleaser release --clean
