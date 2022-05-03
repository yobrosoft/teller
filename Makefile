VERSION := $(shell git describe --tags)
BUILD := $(shell git rev-parse --short HEAD)
NAME := $(shell basename "${PWD}")

GOBASE := $(shell pwd)
GOPATH := $(shell go env GOPATH)
GOBIN := $(GOBASE)/out
GOFILES := $(wildcard cmd/teller/*.go)

LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.BUILD=$(BUILD)"

build:
	@echo "> building ${NAME}..."
	@mkdir -p $(GOBIN)
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) GO11MODULE=on go build $(LDFLAGS) -o $(GOBIN)/$(NAME) $(GOFILES)
