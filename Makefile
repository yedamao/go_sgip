# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test

all: unit build

.PHONY: unit
unit: ## @testing Run the unit tests
	$(GOTEST) -race -coverprofile=coverage.txt -covermode=atomic $(shell go list ./sgip/...)

.PHONY: build
build:
	$(GOBUILD) -o ./bin/receiver ./cmd/receiver
	$(GOBUILD) -o ./bin/transmitter ./cmd/transmitter
	$(GOBUILD) -o ./bin/mockserver ./cmd/mockserver
	$(GOBUILD) -o ./bin/mockclient ./cmd/mockclient

.PHONY: build_linux
build_linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GOBUILD) -o ./bin/receiver ./cmd/receiver
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GOBUILD) -o ./bin/transmitter ./cmd/transmitter
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GOBUILD) -o ./bin/mockserver ./cmd/mockserver
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GOBUILD) -o ./bin/mockclient ./cmd/mockclient

.PHONY: clean
clean:
	rm -rf ./bin coverage.txt
