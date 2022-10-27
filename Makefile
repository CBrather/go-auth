ifneq (,$(wildcard .env))
    include .env
    export
endif

DATE    			?= $(shell date +%FT%T%z)
VERSION 			?= $(shell git describe --tags --always --dirty 2> /dev/null || echo v0)
PKGS				= $(or $(PKG),$(shell env $(GO) list ./...))
BIN					= $(CURDIR)/bin
LOCAL				= $(HOME)/.local
GO					= go

V = 0
Q = $(if $(filter 1,$V),,@)
M = $(shell printf "\033[34;1m▶\033[0m")

# Build
.PHONY: build
build: $(BIN) ; $(info $(M) building release...)
	$(GO) build \
		-ldflags '-X main.Version=$(VERSION) -X main.BuildDate=$(DATE)' \
		-a -installsuffix cgo \
		-tags release \
		-o $(BIN)/main \
		./cmd/main.go;

.PHONY: release
release: $(BIN) ; $(info $(M) building linux release...)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build \
		-ldflags '-X main.Version=$(VERSION) -X main.BuildDate=$(DATE) -s -w -extldflags "-static"' \
		-a -installsuffix cgo \
		-tags release \
		-o $(BIN)/main \
		./cmd/main.go;

# Tools
$(BIN):
	$Q mkdir -p $@

$(GOBIN)/%: ; $(info $(M) installing $(PACKAGE)...)
	$Q $(GO) install $(ARGS) $(PACKAGE)

$(LOCAL)/%: ; $(info $(M) downloading and installing $(URL))
	$Q curl -LO $(URL)
	$Q unzip $(shell basename $(URL)) -d $(LOCAL)
	$Q rm $(shell basename $(URL))

GOLINT = $(GOBIN)/golangci-lint
$(GOBIN)/golangci-lint: PACKAGE=github.com/golangci/golangci-lint/cmd/golangci-lint@v1.47.3

GOCOV = $(GOBIN)/gocov
$(GOBIN)/gocov: PACKAGE=github.com/axw/gocov/gocov@v1.1.0

GOCOVXML = $(GOBIN)/gocov-xml
$(GOBIN)/gocov-xml: PACKAGE=github.com/AlekSi/gocov-xml@v1.0.0

GOJUNIT = $(GOBIN)/go-junit-report
$(GOBIN)/go-junit-report: PACKAGE=github.com/jstemmer/go-junit-report@v1.0.0

# Tests
.PHONY: test tests
test tests: ; $(info $(M) running unit tests…) @ 				## Run unit tests
	$Q $(GO) test -v $(PKGS)

COVERAGE_MODE    = atomic
COVERAGE_PROFILE = $(COVERAGE_DIR)/profile.out
COVERAGE_XML     = $(COVERAGE_DIR)/coverage.xml
COVERAGE_HTML    = $(COVERAGE_DIR)/index.html

.PHONY: test-coverage test-coverage-tools test-integration test-azure
test-coverage-tools: | $(GOCOV) $(GOCOVXML) $(GOJUNIT)
test-coverage: COVERAGE_DIR := $(CURDIR)/test/coverage
test-coverage: test-coverage-tools ; $(info $(M) running integration tests...)
	$Q mkdir -p $(COVERAGE_DIR)
	$Q $(GO) test -v -tags=$(BUILD_TAGS) \
		-coverpkg=$$($(GO) list -f '{{ join .Deps "\n" }}' $(PKGS) | \
					grep '^$(MODULE)/' | \
					tr '\n' ',' | sed 's/,$$//') \
		-covermode=$(COVERAGE_MODE) \
		-coverprofile="$(COVERAGE_PROFILE)" $(PKGS) | tee test/tests.output
	$Q $(GO) tool cover -html=$(COVERAGE_PROFILE) -o $(COVERAGE_HTML)
	$Q $(GOCOV) convert $(COVERAGE_PROFILE) | $(GOCOVXML) > $(COVERAGE_XML)
	$Q cat test/tests.output | $(GOJUNIT) -set-exit-code > test/tests.xml

.PHONY: lint lint-staged lint-no-fix
lint-no-fix: $(GOLINT) ; $(info $(M) running golangci-lint (without fixing issues)…) @ ## Run golangci-lint on codebase (without fixing issues)
	$Q $(GOLINT) run --sort-results
lint-staged: NEW=--new ; $(info $(M) running golangci-lint on staged files...) ## Run golangci-lint on staged source files only
lint-staged: lint
lint: $(GOLINT) ; $(info $(M) running golangci-lint…) @ ## Run golangci-lint on codebase
	$Q $(GOLINT) run --sort-results --fix $(NEW)

.PHONY: clean rm
clean rm: ; $(info $(M) cleaning artifacts...)	@ ## Cleanup everything
	$Q rm -rf $(BIN)
	$Q rm -rf test

.PHONY: help
help:
	$Q grep -hE '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-17s\033[0m %s\n", $$1, $$2}'

.PHONY: version
version:
	$Q echo $(VERSION)
