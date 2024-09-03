# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet
GOMOD=$(GOCMD) mod
BINARY_NAME=oakestra-monitoring-manager
PKG=./...

# Directories
SRC_DIR=./cmd/$(BINARY_NAME)
BINARY_DIR=./bin
BUILD_DIR=./build
CONFIG_DIR=./config
BUILDER_CONFIG_DIR=$(CONFIG_DIR)/opentelemetry-collector-builder
SCRIPTS_DIR=./scripts

# Detect OS and arch
OS := $(shell uname -s)
ARCH := $(shell uname -m)

ifeq ($(OS), Darwin)
	ifeq ($(ARCH), arm64)
		OSARCH=darwin-arm64
	endif
endif

ifeq ($(OS), Linux)
	ifeq ($(ARCH), aarch64)
		OSARCH=linux-aarch64
	endif
endif
# default
OSARCH ?= unknown

# OpenTelemetry Collector Builder
OCB=builder
#OCB=$(SCRIPTS_DIR)/ocb-$(OSARCH)
COLLECTOR_BIN=manager
COLLECTOR_BUILD_DIR=/tmp/oakestra-monitoring-manager
COLLECTOR_CONFIG_DIR=$(CONFIG_DIR)/opentelemetry-collector

# Versioning
VERSION=$(shell git describe --tags --always)
COMMIT=$(shell git rev-parse --short HEAD)
DATE=$(shell date +%Y-%m-%dT%H:%M:%SZ)

# Go build flags
LDFLAGS=-ldflags "-X 'main.Version=$(VERSION)' -X 'main.Commit=$(COMMIT)' -X 'main.Date=$(DATE)'"

# Variables
BINARY_NAME ?= oakestra-monitoring-manager

# Docker
DOCKER=docker
CONTAINER_NAME=monitoring-manager

# Commands
all: help

help:
	@echo "Available commands:"
	@echo "build-docker \t\t build docker image"
	@echo "build \t\t\t build opentelemetry collector"
	@echo "test \t\t\t run tests"
	@echo "fmt \t\t\t run gofmt"
	@echo "vet \t\t\t run go vet"
	@echo "mod-tidy \t\t run go mod tidy"
	@echo "run \t\t\t run the collector in $(COLLECTOR_BUILD_DIR)/$(COLLECTOR_BIN)"

.PHONY: push
push: build-export
	@echo "Pushing bin and build to github..."
	@git add bin
	@git add build
	@git commit -m "Pushed build $(COMMIT) on $(DATE)"
	@git push

.PHONY: build-export
build-export: build clean
	@echo "(Re-)building project for export..."
	@cp -r $(COLLECTOR_BUILD_DIR) $(BUILD_DIR)
	@mkdir -p $(BINARY_DIR)
	@mv $(BUILD_DIR)/$(COLLECTOR_BIN) $(BINARY_DIR)/$(COLLECTOR_BIN)
	@echo "Done"

.PHONY: build-docker
build-docker:
  docker build --progress=plain -t monitoring-manager:latest .

.PHONY: build
build:
	@echo "Building collector... "
	$(OCB) --config=$(BUILDER_CONFIG_DIR)/manifest.yaml --name=$(COLLECTOR_BIN) --output-path=$(COLLECTOR_BUILD_DIR) --skip-strict-versioning

.PHONY: clean
clean:
	@echo "Cleaning up..."
	@rm -r $(BUILD_DIR)
	@rm -r $(BINARY_DIR)

.PHONY: test
test:
	@echo "Running tests..."
	$(GOTEST) -v $(PKG)

.PHONY: fmt
fmt:
	@echo "Formatting code..."
	$(GOFMT) $(PKG)

.PHONY: vet
vet:
	@echo "Vetting code..."
	$(GOVET) $(PKG)

.PHONY: mod-tidy
mod-tidy: 
	@echo "Tidying up modules..."
	$(GOMOD) tidy

.PHONY: run
run: build
	@echo "Running collector..."
	$(COLLECTOR_BUILD_DIR)/$(COLLECTOR_BIN) --config=$(COLLECTOR_CONFIG_DIR)/opentelemetry-config.yaml
