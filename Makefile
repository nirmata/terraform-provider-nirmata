.DEFAULT_GOAL      := build
OS                 := $(shell go env GOOS)
ARCH               := $(shell go env GOARCH)
VERSION        	   := $(shell git describe)
PLUGINS            ?=${HOME}/.terraform.d/plugins
PLUGIN_PATH        ?=local/nirmata/nirmata/${VERSION}/${OS}_${ARCH}
PLUGIN_NAME        := terraform-provider-nirmata_${VERSION}
DIST_PATH          := dist/${OS}_${ARCH}
GO_PACKAGES        := $(shell go list ./... | grep -v /vendor/)
GO_FILES           := $(shell find . -type f -name '*.go')
GO                 ?= go

.PHONY: all
all: test build

.PHONY: test
test: test-all

.PHONY: test-all
test-all:
	@TF_ACC=1 $(GO) test -v -race $(GO_PACKAGES)

${DIST_PATH}/${PLUGIN_NAME}: ${GO_FILES}
	mkdir -p $(DIST_PATH); \
	$(GO) build -o $(DIST_PATH)/${PLUGIN_NAME}

.PHONY: build
build: ${DIST_PATH}/${PLUGIN_NAME}

.PHONY: install
install: clean build
	mkdir -p $(PLUGIN_PATH); \
	rm -rf ${PLUGINS}/$(PLUGIN_PATH)/${PLUGIN_NAME}; \
	install -m 0755 "$(DIST_PATH)/${PLUGIN_NAME}" "${PLUGINS}/$(PLUGIN_PATH)/${PLUGIN_NAME}"

.PHONY: clean
clean:
	rm -rf ${DIST_PATH}/*
