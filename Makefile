BIN_NAME=vqs

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BUILD_DIR=bin

build:
	@CGO_ENABLED=0 $(GOBUILD) -o $(BUILD_DIR)/$(BIN_NAME)

run:
	go run main.go
