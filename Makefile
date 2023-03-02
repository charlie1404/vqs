BIN_NAME=vqs

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BUILD_DIR=bin
MODULE=$(shell go list -m)

VFLAG=-X '$(MODULE)/cmd.CURRENT_VERSION=1.2.1'

.PHONY: build run

build:
	$(GOBUILD) -ldflags "$(VFLAG)" -o $(BUILD_DIR)/$(BIN_NAME)

run:
	go run main.go