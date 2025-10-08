.PHONY := all build test lint clean cover

GO        ?= go
BIN_DIR   ?= bin
BINARY    ?= tms_ssh_emulator
MAIN_PKG  ?= ./src
MODULE    := $(shell $(GO) list -m)
ALL_PKGS  := $(shell $(GO) list ./...)
# Exclude the main package from unit test coverage; it's an executable entrypoint.
PKGS      ?= $(filter-out $(MODULE)/src,$(ALL_PKGS))

all: build

build:
	@mkdir -p $(BIN_DIR)
	$(GO) build -trimpath -ldflags="-s -w" -o $(BIN_DIR)/$(BINARY) $(MAIN_PKG)

test:
	$(GO) test $(PKGS) -race -coverprofile=coverage.out -covermode=atomic -count=1

lint:
	$(GO) vet $(PKGS)

clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BIN_DIR) dist out
	rm -f coverage.*.out coverage.out
	$(GO) clean -testcache

cover: test
	@mkdir -p out
	$(GO) tool cover -func=coverage.out
	$(GO) tool cover -html=coverage.out -o out/coverage.html
	@echo "HTML coverage report: out/coverage.html"
