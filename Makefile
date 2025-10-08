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

# --- Semantic version tagging -------------------------------------------------
.PHONY: tag/patch tag/minor tag/major

define _SEMVER_TAG
@set -e; \
latest=$$(git tag -l 'v[0-9]*.[0-9]*.[0-9]*' --sort=-v:refname | head -n1); \
if [ -z "$$latest" ]; then \
  base=v0.0.0; \
  echo "No tags found. Creating initial $$base"; \
  git tag -a "$$base" -m "chore(tag): initialize $$base"; \
  git push origin "$$base"; \
  latest="$$base"; \
fi; \
nums=$${latest#v}; \
IFS=.; set -- $$nums; \
major=$$1; minor=$$2; patch=$$3; \
case "$(1)" in \
  patch) new_major=$$major; new_minor=$$minor; new_patch=$$((patch+1));; \
  minor) new_major=$$major; new_minor=$$((minor+1)); new_patch=0;; \
  major) new_major=$$((major+1)); new_minor=0; new_patch=0;; \
  *) echo "unknown bump kind: $(1)" >&2; exit 2;; \
esac; \
new_tag=v$$new_major.$$new_minor.$$new_patch; \
echo "Latest: $$latest -> New: $$new_tag"; \
git tag -a "$$new_tag" -m "chore(tag): release $$new_tag"; \
git push origin "$$new_tag"
endef

tag/patch:
	$(call _SEMVER_TAG,patch)

tag/minor:
	$(call _SEMVER_TAG,minor)

tag/major:
	$(call _SEMVER_TAG,major)

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
