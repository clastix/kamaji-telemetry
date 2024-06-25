DOCKER_REPO ?= docker.io/clastix/kamaji-telemetry

GOLANGCI_LINT 			= $(shell pwd)/bin/golangci-lint
GOLANGCILINT_VERSION	?= v1.59.1

KO 			= $(shell pwd)/bin/ko
KO_TAGS 	?= latest
KO_VERSION	?= v0.14.1
KO_CACHE 	?= /tmp/ko-cache
KO_LOCAL    ?= true

# go-install-tool will 'go install' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-install-tool
@[ -f $(1) ] || { \
set -e ;\
echo "Installing $(2)" ;\
GOBIN=$(PROJECT_DIR)/bin go install $(2) ;\
}
endef

golangci-lint: ## Download golangci-lint locally if necessary.
	$(call go-install-tool,$(GOLANGCI_LINT),github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCILINT_VERSION))

lint: golangci-lint ## Linting the code according to the styling guide.
	$(GOLANGCI_LINT) run -c .golangci.yml

.PHONY: ko
ko:
	$(call go-install-tool,$(KO),github.com/google/ko@$(KO_VERSION))

.PHONY: ko-build
ko-build: ko
	KOCACHE=$(KO_CACHE) KO_DOCKER_REPO=$(DOCKER_REPO) \
	$(KO) build ./cmd/server --local=$(KO_LOCAL) --bare --tags=$(KO_TAGS)
