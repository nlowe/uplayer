GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

ARCH=amd64
DIST=".$(notdir $(patsubst %/,%,$(dir $($(abspath $(lastword $(MAKEFILE_LIST)))))))/_dist"
BINARY_NAME=uplayer

all: build
build: deps build-linux build-windows

build-linux: build-prepare
	GOOS=linux GOARCH="$(ARCH)" $(GOBUILD) -o "$(DIST)/linux/$(BINARY_NAME)" -v

build-windows: build-prepare
	GOOS=windows GOARCH="$(ARCH)" $(GOBUILD) -o "$(DIST)/windows/$(BINARY_NAME).exe" -v

build-prepare:
	rm -rf "$(DIST)"
	mkdir -p "$(DIST)"
	mkdir -p "$(DIST)/linux"
	mkdir -p "$(DIST)/windows"

clean: 
		$(GOCLEAN)
		rm -rf "$(DIST)"

run:
		$(GOBUILD) -o $(BINARY_NAME) -v ./...
		./$(BINARY_NAME)

deps:
		$(GOGET) .
