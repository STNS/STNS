TEST ?= $(shell go list ./... | grep -v vendor)
VERSION = $(shell cat version)
REVISION = $(shell git describe --always)

INFO_COLOR=\033[1;34m
RESET=\033[0m
BOLD=\033[1m

default: build
ci: depsdev test lint integration ## Run test and more...

deps: ## Install dependencies
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Installing Dependencies$(RESET)"
	go get -u golang.org/x/vgo/...
	vgo install

depsdev: deps ## Installing dependencies for development
	go get github.com/golang/lint/golint
	go get github.com/pierrre/gotestcover
	go get -u github.com/tcnksm/ghr
	go get github.com/mitchellh/gox

test: ## Run test
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Testing$(RESET)"
	vgo test -v $(TEST) -timeout=30s -parallel=4
	vgo test -race $(TEST)

lint: ## Exec golint
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Linting$(RESET)"
	golint -min_confidence 1.1 -set_exit_status $(TEST)

server: ## Run server
	vgo run github.com/STNS/STNS --logfile ./stns.log --pidfile ./stns.pid --config ./stns/test.toml server

ghr: ## Upload to Github releases without token check
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Releasing for Github$(RESET)"
	ghr -u stns v$(VERSION)-$(REVISION) pkg

dist: build ## Upload to Github releases
	@test -z $(GITHUB_TOKEN) || test -z $(GITHUB_API) || $(MAKE) ghr

integration: ## Run integration test after Server wakeup
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Integration Testing$(RESET)"
	./misc/server start
	vgo test $(VERBOSE) -integration $(TEST) $(TEST_OPTIONS)
	./misc/server stop

build: ## Build server
	GOOS=linux GOARCH=amd64 vgo build -o misc/stns

docker:
	docker build -t nss_develop .
	docker run --cap-add=SYS_PTRACE --security-opt="seccomp=unconfined" -v `pwd`:/go/src/github.com/STNS/STNS -w /go/src/github.com/STNS/STNS -it nss_develop /bin/bash

.PHONY: default dist test deps docker
