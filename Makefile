TEST ?= $(shell go list ./... | grep -v vendor)
VERSION = $(shell cat version)
REVISION = $(shell git describe --always)

INFO_COLOR=\033[1;34m
RESET=\033[0m
BOLD=\033[1m

default: build
ci: depsdev test vet lint ## Run test and more...

deps: ## Install dependencies
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Installing Dependencies$(RESET)"
	go get -u golang.org/x/vgo/...
	vgo install

depsdev: deps ## Installing dependencies for development
	go get github.com/golang/lint/golint
	go get github.com/codegangsta/gin
	go get github.com/pierrre/gotestcover
	go get -u github.com/tcnksm/ghr
	go get github.com/mitchellh/gox

test: ## Run test
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Testing$(RESET)"
	vgo test -v $(TEST) -timeout=30s -parallel=4
	vgo test -race $(TEST)

vet: ## Exec go vet
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Vetting$(RESET)"
	vgo vet $(TEST)

lint: ## Exec golint
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Linting$(RESET)"
	golint -min_confidence 1.1 -set_exit_status $(TEST)

server: ## Run server with gin
	gin -appPort 8050

ghr: ## Upload to Github releases without token check
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Releasing for Github$(RESET)"
	ghr -u stns v$(VERSION)-$(REVISION) pkg

dist: build ## Upload to Github releases
	@test -z $(GITHUB_TOKEN) || test -z $(GITHUB_API) || $(MAKE) ghr


.PHONY: default dist test deps 
