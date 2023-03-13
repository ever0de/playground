PACKAGES=$(shell go list ./...)

all: lint

lint: 
	@echo "--> Running linter"
	@golangci-lint run
	@go mod verify
.PHONY: lint
