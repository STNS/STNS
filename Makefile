VERSION = $(shell git describe --tags --abbrev=0|sed -e 's/v//g')
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
SOURCES=Makefile v2/go.mod v2/go.sum version v2 package/
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
GOOS=linux
GOARCH=amd64
GO=GO111MODULE=on GOOS=$(GOOS) GOARCH=$(GOARCH) go

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
	cd $(PACKAGE_DIR) && $(GO) test $(TEST_LIST) -race

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

source_for_rpm: ## Create source for RPM
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Distributing$(RESET)"
	rm -rf tmp.$(DIST) stns-v2-$(VERSION).tar.gz
	mkdir -p tmp.$(DIST)/stns-v2-$(VERSION)
	cp -r $(SOURCES) tmp.$(DIST)/stns-v2-$(VERSION)
	mkdir -p tmp.$(DIST)/stns-v2-$(VERSION)/tmp/bin
	cp -r tmp/bin/* tmp.$(DIST)/stns-v2-$(VERSION)/tmp/bin
	cd tmp.$(DIST) && \
		tar cf stns-v2-$(VERSION).tar stns-v2-$(VERSION) && \
		gzip -9 stns-v2-$(VERSION).tar
	cp tmp.$(DIST)/stns-v2-$(VERSION).tar.gz ./builds
	rm -rf tmp.$(DIST)

rpm: source_for_rpm ## Packaging for RPM
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Packaging for RPM$(RESET)"
	cp builds/stns-v2-$(VERSION).tar.gz /root/rpmbuild/SOURCES
	spectool -g -R rpm/stns.spec
	rpmbuild -ba rpm/stns.spec
	cp /root/rpmbuild/RPMS/*/*.rpm /go/src/github.com/STNS/STNS/builds

SUPPORTOS=centos7 almalinux9 ubuntu20 ubuntu22 debian10 debian11
pkg: build ## Create some distribution packages
	rm -rf builds && mkdir builds
	for i in $(SUPPORTOS); do \
	  docker-compose build $$i; \
	  docker-compose run -v `pwd`:/go/src/github.com/STNS/STNS -v ~/pkg:/go/pkg --rm $$i; \
	done

source_for_deb: ## Create source for DEB
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Distributing$(RESET)"
	rm -rf tmp.$(DIST) stns-v2-$(VERSION).orig.tar.gz
	mkdir -p tmp.$(DIST)/stns-v2-$(VERSION)
	cp -r $(SOURCES) tmp.$(DIST)/stns-v2-$(VERSION)
	mkdir -p tmp.$(DIST)/stns-v2-$(VERSION)/tmp/bin
	cp -r tmp/bin/* tmp.$(DIST)/stns-v2-$(VERSION)/tmp/bin
	cd tmp.$(DIST) && \
	tar zcf stns-v2-$(VERSION).tar.gz stns-v2-$(VERSION)
	mv tmp.$(DIST)/stns-v2-$(VERSION).tar.gz tmp.$(DIST)/stns-v2-$(VERSION).orig.tar.gz

deb: source_for_deb ## Packaging for DEB
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Packaging for DEB$(RESET)"
	cd tmp.$(DIST) && \
		tar xf stns-v2-$(VERSION).orig.tar.gz && \
		cd stns-v2-$(VERSION) && \
		dh_make --single --createorig -y && \
		rm -rf debian/*.ex debian/*.EX debian/README.Debian && \
		cp -r $(GOPATH)/src/github.com/STNS/STNS/debian/* debian/ && \
		sed -i -e 's/xenial/$(DIST)/g' debian/changelog && \
		debuild -uc -us
	cd tmp.$(DIST) && \
		find . -name "*.deb" | sed -e 's/\(\(.*stns-v2.*\).deb\)/mv \1 \2.$(DIST).deb/g' | sh && \
		cp *.deb $(GOPATH)/src/github.com/STNS/STNS/builds
	rm -rf tmp.$(DIST)

github_release: ## Create some distribution packages
	ghr -u STNS --replace v$(VERSION) builds/

generate:
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Generate From ERB$(RESET)"
	ruby model/make_backends.rb

.PHONY: default test docker rpm source_for_rpm pkg source_for_deb deb server
