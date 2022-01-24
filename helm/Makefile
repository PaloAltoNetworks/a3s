export PROJECT_VERSION ?= v0.0.0-dev

default: lint charts

lint:
	helm lint ./a3s --values ./test/values.yaml

charts: lint
	helm package a3s --version $(PROJECT_VERSION) -d .
	helm push --force ./a3s-$(PROJECT_VERSION).tgz local;
	helm repo update