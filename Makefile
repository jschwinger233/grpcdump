.PHONY: deps binary build

REPO_PATH := github.com/jschwinger233/grpcdump
REVISION := $(shell git rev-parse HEAD || unknown)
BUILTAT := $(shell date +%Y-%m-%dT%H:%M:%S)
VERSION := $(shell git describe --tags $(shell git rev-list --tags --max-count=1))
GO_LDFLAGS ?= -X $(REPO_PATH)/version.REVISION=$(REVISION) \
			  -X $(REPO_PATH)/version.BUILTAT=$(BUILTAT) \
			  -X $(REPO_PATH)/version.VERSION=$(VERSION)

build: deps binary

deps:
	env GO111MODULE=on go mod download
	env GO111MODULE=on go mod vendor

binary:
	go build -ldflags "$(GO_LDFLAGS)"
