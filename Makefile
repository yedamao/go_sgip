# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test

all: unit build

.PHONY: unit
unit: ## @testing Run the unit tests
	$(GOTEST) $(shell go list ./sgip/...)

.PHONY: build
build: clean
	$(GOBUILD) ./cmd/receiver
	$(GOBUILD) ./cmd/transmitter

.PHONY: build_linux
build_linux: clean
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GOBUILD) -o receiver_linux  ./cmd/receiver
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GOBUILD) -o transmitter_linx ./cmd/transmitter

.PHONY: clean
clean:
	$(GOCLEAN)
	rm -rf receiver transmitter receiver_linux transmitter_linx
