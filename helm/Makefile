export PROJECT_VERSION ?= v0.0.0-dev

all: lint charts

lint:
	helm lint ./a3s --values ./test/values.yaml

charts: lint
	helm package a3s --version $(PROJECT_VERSION) -d .