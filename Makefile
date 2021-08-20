MAKEFLAGS += --warn-undefined-variables
SHELL := /bin/bash
ifeq ($(word 1,$(subst ., ,$(MAKE_VERSION))),4)
.SHELLFLAGS := -eu -o pipefail -c
endif
.DEFAULT_GOAL := bin/terraform-provider-dotfiles
.ONESHELL:

GIT_REF := $(shell git describe --match="" --always --dirty=+)
GIT_TAG := $(shell git name-rev --tags --name-only $(GIT_REF))
PACKAGE := $(shell go list)

.PHONY: help
help:  ## Show this help
	@grep '.*:.*##' Makefile | grep -v grep  | sort | sed 's/:.* ##/:/g' | column -t -s:

.PHONY: test
test:  ## Run go test
	go test -v ./...

bin/terraform-provider-dotfiles:  ## Build the application binary for target OS, for example bin/terraform-provider-dotfiles.linux
	go build -o $@ -ldflags "-X $(PACKAGE)/version=$(GIT_TAG)+$(GIT_REF)" main.go

.PHONY: install
install: bin/terraform-provider-dotfiles
	mkdir -p terraform.d/plugins/github.com/mhristof/dotfiles/0.1.0/darwin_amd64
	mv $< terraform.d/plugins/github.com/mhristof/dotfiles/0.1.0/darwin_amd64

