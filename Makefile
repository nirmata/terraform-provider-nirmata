.DEFAULT_GOAL      := build
OS                 := $(shell go env GOOS)
ARCH               := $(shell go env GOARCH)
VERSION        	   := $(shell git describe)

PLUGIN_NAME        := terraform-provider-nirmata_${VERSION}
DIST_PATH          := dist/${OS}_${ARCH}

.PHONY: all
all: build

.PHONY: build
build:
	mkdir -p $(DIST_PATH); \
	go build -o $(DIST_PATH)/${PLUGIN_NAME}

.PHONY: clean
clean:
	rm -rf ${DIST_PATH}/*
