MAKEFLAGS += --warn-undefined-variables
SHELL := /bin/bash
ifeq ($(word 1,$(subst ., ,$(MAKE_VERSION))),4)
.SHELLFLAGS := -eu -o pipefail -c
endif
.DEFAULT_GOAL := bin/terraform-provider-dotfiles.darwin
.ONESHELL:

PACKAGE := $(shell go list)
GIT_REF := $(shell git describe --match="" --always --dirty=+)
GIT_TAG := $(shell git name-rev --tags --name-only $(GIT_REF))

.PHONY: help
help:  ## Show this help
	@grep '.*:.*##' Makefile | grep -v grep  | sort | sed 's/:.* ##/:/g' | column -t -s:

.PHONY: init
init: .terraform ## Force run 'terraform init'

.terraform:  ## 
	terraform init

.PHONY: plan
plan: terraform.tfplan ## Runs 'terraform plan'

terraform.tfplan: $(shell find ./ -name '*.tf') .terraform ## Creates terraform.tfplan if required
	terraform plan -out $@

.PHONY: apply
apply: terraform.tfstate ## Run 'terraform apply'

terraform.tfstate: terraform.tfplan ## Run 'terraform apply' if required'
	terraform apply terraform.tfplan

.PHONY: force
force:  ## Forcefully update terraform state
	touch *.tf && make terraform.tfstate

.PHONY: destroy
destroy:  ## Run 'terraform destroy'
	terraform destroy -auto-approve

.PHONY: clean
clean: destroy ## Clean the repository resources
	rm -rf terraform.tf{state,plan} .terraform terraform.state.d
	rm bin/*

.PHONY: test
test:  ## Run go test
	go test -v ./...

bin/terraform-provider-dotfiles.darwin:  ## Build the application binary for current OS

bin/terraform-provider-dotfiles.%:  ## Build the application binary for target OS, for example bin/terraform-provider-dotfiles.linux
	GOOS=$* go build -o $@ -ldflags "-X $(PACKAGE)/version=$(GIT_TAG)+$(GIT_REF)" main.go

.PHONY: install
install: bin/terraform-provider-dotfiles.darwin ## Install the binary
	cp $< ~/bin/terraform-provider-dotfiles

unknown flag: --go
