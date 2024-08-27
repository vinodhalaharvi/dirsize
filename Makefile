# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOVET=$(GOCMD) vet
BINARY_NAME=dirsize
BINARY_UNIX=$(BINARY_NAME)_unix

# Linter
GOLINT=golangci-lint

# Define the installation directory
INSTALL_DIR=$(HOME)/go/bin

all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)

deps:
	$(GOGET) -v ./...

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v

# Linting
lint:
	$(GOLINT) run

# Vetting
vet:
	$(GOVET) ./...

# All quality checks
quality: lint vet test

# Install tools
install-tools:
	@mkdir -p $(INSTALL_DIR)
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(INSTALL_DIR) v1.55.2
	@echo "golangci-lint installed successfully in $(INSTALL_DIR)"
	@echo "Make sure $(INSTALL_DIR) is in your PATH"

.PHONY: all build test clean run deps build-linux lint vet quality install-tools