export DOCKER_TAG ?= v0.0.0-dev

default: lint charts

lint:
	helm lint ./a3s --values ./test/values.yaml

charts: lint
	helm package a3s --version $(DOCKER_TAG) -d .
	helm push --force ./a3s-$(DOCKER_TAG).tgz local;
	helm repo update