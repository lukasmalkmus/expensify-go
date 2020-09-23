# TOOLCHAIN
GO		:= CGO_ENABLED=0 GOBIN=$(CURDIR)/bin go
GOFMT	:= $(GO)fmt

# ENVIRONMENT
VERBOSE 	=
GOPATH		:= $(GOPATH)
GOOS		?= $(shell echo $(shell uname -s) | tr A-Z a-z)
GOARCH		?= amd64
MOD_NAME	:= github.com/lukasmalkmus/expensify-go

# TOOLS
GOLANGCI_LINT	:= bin/golangci-lint
GOTESTSUM		:= bin/gotestsum

# MISC
COVERPROFILE	:= coverage.out

# FLAGS
GOTESTSUM_FLAGS	:= --jsonfile tests.json --junitfile junit.xml
GO_TEST_FLAGS 	:= -race -coverprofile=$(COVERPROFILE)

# DEPENDENCIES
GOMODDEPS = go.mod go.sum

# Enable verbose test output if explicitly set.
ifdef VERBOSE
	GOTESTSUM_FLAGS	+= --format=standard-verbose
endif

# FUNCS
# func go-list-pkg-sources(package)
go-list-pkg-sources = $(GO) list $(GOFLAGS) -f '{{ range $$index, $$filename := .GoFiles }}{{ $$.Dir }}/{{ $$filename }} {{end}}' $(1)
# func go-pkg-sourcefiles(package)
go-pkg-sourcefiles = $(shell $(call go-list-pkg-sources,$(strip $1)))

.PHONY: all
all: dep fmt lint test ## Run dep, fmt, lint and test.

.PHONY: clean
clean: ## Remove test artifacts.
	@echo ">> cleaning up artifacts"
	@rm -rf $(COVERPROFILE)

.PHONY: cover
cover: $(COVERPROFILE) ## Calculate the code coverage score.
	@echo ">> calculating code coverage"
	@$(GO) tool cover -func=$(COVERPROFILE)

.PHONY: dep-clean
dep-clean: ## Remove obsolete dependencies.
	@echo ">> cleaning dependencies"
	@$(GO) mod tidy

.PHONY: dep-upgrade
dep-upgrade: ## Upgrade all direct dependencies to their latest version.
	@echo ">> upgrading dependencies"
	@$(GO) get $(shell $(GO) list -f '{{if not (or .Main .Indirect)}}{{.Path}}{{end}}' -m all)
	@make dep

.PHONY: dep
dep: dep-clean dep.stamp ## Install dependencies.

dep.stamp: $(GOMODDEPS)
	@echo ">> installing dependencies"
	@$(GO) mod download
	@$(GO) mod verify
	@touch $@

.PHONY: fmt
fmt: ## Format and simplify the source code using `gofmt`.
	@echo ">> formatting code"
	@! $(GOFMT) -s -w $(shell find . -path -prune -o -name '*.go' -print) | grep '^'

.PHONY: lint
lint: $(GOLANGCI_LINT) ## Lint the source code.
	@echo ">> linting code"
	@$(GOLANGCI_LINT) run

.PHONY: test
test: $(GOTESTSUM) ## Run all tests. Run with VERBOSE=1 to get verbose test output ('-v' flag).
	@echo ">> running tests"
	@$(GOTESTSUM) $(GOTESTSUM_FLAGS) -- $(GO_TEST_FLAGS) ./...

.PHONY: tools
tools: $(GOLANGCI_LINT) $(GOTESTSUM) ## Install all tools into the projects local $GOBIN directory.

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# TEST TARGETS

$(COVERPROFILE):
	@make test

# TOOLS

$(GOLANGCI_LINT): dep.stamp $(call go-pkg-sourcefiles, github.com/golangci/golangci-lint/cmd/golangci-lint)
	@echo ">> installing golangci-lint"
	@$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint

$(GOTESTSUM): dep.stamp $(call go-pkg-sourcefiles, gotest.tools/gotestsum)
	@echo ">> installing gotestsum"
	@$(GO) install gotest.tools/gotestsum
