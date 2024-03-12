VERSION = $(shell git tag | sed 's/v//g' |sort --version-sort | tail -n1)
REVISION = $(shell git describe --always)
TEST_LIST = $(shell cd v2 && go list ./...)

INFO_COLOR=\033[1;34m
RESET=\033[0m
BOLD=\033[1m

PACKAGE_DIR="v2"
DIST ?= unknown
PREFIX=/usr
BINDIR=$(PREFIX)/sbin
MODDIR ?= $(PREFIX)/local/stns/modules.d
SOURCES=Makefile v2/go.mod v2/go.sum v2 package/
ETCD_VER=3.3.10
REDIS_VER=5.0.4
BUILD:=$(shell pwd)/tmp/bin
MIDDLEWARE:=$(shell pwd)/tmp/middleware
UNAME_S := $(shell uname -s)

REVISION=$(shell git describe --always)
GOVERSION=$(shell go version)
BUILDDATE=$(shell date '+%Y/%m/%d %H:%M:%S %Z')
STNS_PROTOCOL ?= "http"
GOPATH ?= /go
GO=GO111MODULE=on go

ME=$(shell whoami)
default: build

ci: depsdev test lint integration ## Run test and more...

etcd:
	echo $(UNAME_S)
	mkdir -p $(MIDDLEWARE)
ifeq ($(UNAME_S),Linux)
	test -e $(MIDDLEWARE)/etcd-v$(ETCD_VER)-linux-amd64/etcd || curl -L  https://github.com/coreos/etcd/releases/download/v$(ETCD_VER)/etcd-v$(ETCD_VER)-linux-amd64.tar.gz -o $(MIDDLEWARE)/etcd-v$(ETCD_VER)-linux-amd64.tar.gz
	test -e $(MIDDLEWARE)/etcd-v$(ETCD_VER)-linux-amd64/etcd || (cd $(MIDDLEWARE) && tar xzf etcd-v$(ETCD_VER)-linux-amd64.tar.gz)
	ps -aux |grep etcd |grep -q -v grep || $(MIDDLEWARE)/etcd-v$(ETCD_VER)-linux-amd64/etcd &
endif
ifeq ($(UNAME_S),Darwin)
	brew services start etcd
endif

redis:
	echo $(UNAME_S)
	mkdir -p $(MIDDLEWARE)
ifeq ($(UNAME_S),Linux)
	test -e $(MIDDLEWARE)/redis-$(REDIS_VER).tar.gz || curl -L  http://download.redis.io/releases/redis-$(REDIS_VER).tar.gz -o $(MIDDLEWARE)//redis-$(REDIS_VER).tar.gz
	test -d $(MIDDLEWARE)/redis-$(REDIS_VER) || (cd $(MIDDLEWARE) && tar xzf redis-$(REDIS_VER).tar.gz)
	test -e $(MIDDLEWARE)/redis-$(REDIS_VER)/src/redis-server || (cd $(MIDDLEWARE)/redis-$(REDIS_VER) && make)
	ps -aux |grep redis |grep -q -v grep || $(MIDDLEWARE)/redis-$(REDIS_VER)/src/redis-server &
endif
ifeq ($(UNAME_S),Darwin)
	brew services start redis
endif

depsdev: ## Installing dependencies for development
	which staticcheck > /dev/null || $(GO) install honnef.co/go/tools/cmd/staticcheck@latest
	which ghr > /dev/null || $(GO) install github.com/tcnksm/ghr@latest
	which goimports > /dev/null ||$(GO) install golang.org/x/tools/cmd/goimports@latest
	which git-chglog > /dev/null ||$(GO) install github.com/git-chglog/git-chglog/cmd/git-chglog@latest
	cd $(PACKAGE_DIR) && $(GO) mod tidy

changelog:
	git-chglog -o CHANGELOG.md

test: ## Run test
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Testing$(RESET) (require: etcd,redis)"
	cd $(PACKAGE_DIR) && $(GO) test $(TEST_LIST) -v -timeout=30s -parallel=4
	cd $(PACKAGE_DIR) && CGO_ENABLE=1 go test $(TEST_LIST) -race

lint: ## Exec golint
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Linting$(RESET)"
	cd $(PACKAGE_DIR) && $(GOPATH)/bin/staticcheck ./...

server: ## Run server
	cd $(PACKAGE_DIR) && $(GO) run github.com/STNS/STNS/v2 --listen 127.0.0.1:1104 --pidfile ./stns.pid --config ./stns/integration.toml --protocol $(STNS_PROTOCOL) server

integration: integration_http integration_ldap ## Run integration test after Server wakeup

integration_http: ## Run integration test after Server wakeup
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Integration HTTP Testing$(RESET)"
	./misc/server start -http
	cd $(PACKAGE_DIR) && $(GO) test $(VERBOSE) -integration-http $(TEST_OPTIONS)
	./misc/server stop || true

integration_ldap: ## Run integration test after Server wakeup
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Integration LDAP Testing$(RESET)"
	./misc/server start -ldap
	cd $(PACKAGE_DIR) && $(GO) test $(VERBOSE) -integration-ldap $(TEST_OPTIONS)
	./misc/server stop || true

build: ## Build server
	git config --global --add safe.directory $(GOPATH)/src/github.com/STNS/STNS
	cd $(PACKAGE_DIR) && $(GO) build -buildvcs=false -ldflags "-X main.version=$(VERSION) -X main.revision=$(REVISION) -X \"main.goversion=$(GOVERSION)\" -X \"main.builddate=$(BUILDDATE)\" -X \"main.builduser=$(ME)\"" -o $(BUILD)/stns
	cd $(PACKAGE_DIR) && $(GO) build -buildvcs=false -buildmode=plugin -o $(BUILD)/mod_stns_etcd.so modules/etcd.go modules/module.go
	cd $(PACKAGE_DIR) && $(GO) build -buildvcs=false -buildmode=plugin -o $(BUILD)/mod_stns_dynamodb.so modules/dynamodb.go modules/module.go

install: ## Install
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Installing as Server$(RESET)"
	cp $(BUILD)/stns $(BINDIR)/stns
	mkdir -p $(MODDIR)/
	cp $(BUILD)/*so $(MODDIR)/

build_image:
	docker build -t stns/stns:$(VERSION) -t stns/stns:latest -t ghcr.io/stns/stns:$(VERSION) -t ghcr.io/stns/stns:latest .

push_image:
	docker push stns/stns:$(VERSION)
	docker push stns/stns:latest
	docker push ghcr.io/stns/stns:$(VERSION)
	docker push ghcr.io/stns/stns:latest

generate:
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Generate From ERB$(RESET)"
	ruby model/make_backends.rb

.PHONY: release
## release: release nke (tagging and exec goreleaser)
release:
	curl -sfL https://goreleaser.com/static/run | bash

.PHONY: default test docker server
