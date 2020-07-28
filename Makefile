TEST ?= $(shell go list ./... | grep -v -e vendor -e keys -e tmp)
VERSION = $(shell cat version)
REVISION = $(shell git describe --always)

GO ?= GO111MODULE=on go
INFO_COLOR=\033[1;34m
RESET=\033[0m
BOLD=\033[1m

PACKAGE_DIR="v2"
DIST ?= unknown
PREFIX=/usr
BINDIR=$(PREFIX)/sbin
MODDIR ?= $(PREFIX)/local/stns/modules.d
SOURCES=Makefile go.mod go.sum version model api middleware modules stns server stns.go package/
DISTS=centos7 centos6 ubuntu16
ETCD_VER=3.3.10
REDIS_VER=5.0.4
BUILD:=$(shell pwd)/tmp/bin
MIDDLEWARE:=$(shell pwd)/tmp/middleware
UNAME_S := $(shell uname -s)

REVISION=$(shell git describe --always)
GOVERSION=$(shell go version)
BUILDDATE=$(shell date '+%Y/%m/%d %H:%M:%S %Z')
STNS_PROTOCOL ?= "http"

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
	$(GO) get -u golang.org/x/lint/golint
	$(GO) get -u github.com/tcnksm/ghr
	$(GO) get -u golang.org/x/tools/cmd/goimports
	$(GO) get -u github.com/git-chglog/git-chglog/cmd/git-chglog
	$(GO) get -u github.com/ugorji/go/codec@none

changelog:
	git-chglog -o CHANGELOG.md

test: ## Run test
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Testing$(RESET) (require: etcd,redis)"
	cd $(PACKAGE_DIR) && $(GO) test -v $(TEST) -timeout=30s -parallel=4
	cd $(PACKAGE_DIR) && $(GO) test -race $(TEST)

lint: ## Exec golint
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Linting$(RESET)"
	cd $(PACKAGE_DIR) && golint -min_confidence 1.1 -set_exit_status $(TEST)

server: ## Run server
	cd $(PACKAGE_DIR) && $(GO) run github.com/STNS/STNS/v2 --listen 127.0.0.1:1104 --pidfile ./stns.pid --config ./stns/integration.toml --protocol $(STNS_PROTOCOL) server

integration: integration_http integration_ldap ## Run integration test after Server wakeup

integration_http: ## Run integration test after Server wakeup
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Integration HTTP Testing$(RESET)"
	./misc/server start -http
	cd $(PACKAGE_DIR) && $(GO) test $(VERBOSE) -integration-http $(TEST) $(TEST_OPTIONS)
	./misc/server stop || true

integration_ldap: ## Run integration test after Server wakeup
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Integration LDAP Testing$(RESET)"
	./misc/server start -ldap
	cd $(PACKAGE_DIR) && $(GO) test $(VERBOSE) -integration-ldap $(TEST) $(TEST_OPTIONS)
	./misc/server stop || true

build: ## Build server
	cd $(PACKAGE_DIR) && $(GO) build -ldflags "-X main.version=$(VERSION) -X main.revision=$(REVISION) -X \"main.goversion=$(GOVERSION)\" -X \"main.builddate=$(BUILDDATE)\" -X \"main.builduser=$(ME)\"" -o $(BUILD)/stns
	cd $(PACKAGE_DIR) && $(GO) build -buildmode=plugin -o $(BUILD)/mod_stns_etcd.so modules/etcd.go modules/module.go
	cd $(PACKAGE_DIR) && $(GO) build -buildmode=plugin -o $(BUILD)/mod_stns_dynamodb.so modules/dynamodb.go modules/module.go

install: build ## Install
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Installing as Server$(RESET)"
	cp $(BUILD)/stns $(BINDIR)/stns
	mkdir -p $(MODDIR)/
	cp $(BUILD)/*so $(MODDIR)/

docker:
	docker build -t stns_develop .
	docker run --cap-add=SYS_PTRACE --security-opt="seccomp=unconfined" -v $(GOPATH):/go/ -v $(GOPATH)/pkg/mod/cache:/go/pkg/mod/cache -w /go/src/github.com/STNS/STNS -it stns_develop /bin/bash

source_for_rpm: ## Create source for RPM
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Distributing$(RESET)"
	rm -rf tmp.$(DIST) stns-v2-$(VERSION).tar.gz
	mkdir -p tmp.$(DIST)/stns-v2-$(VERSION)
	cp -r $(SOURCES) tmp.$(DIST)/stns-v2-$(VERSION)
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

pkg: ## Create some distribution packages
	rm -rf builds && mkdir builds
	docker-compose run -v `pwd`:/go/src/github.com/STNS/STNS -v ~/pkg:/go/pkg --rm centos6
	docker-compose run -v `pwd`:/go/src/github.com/STNS/STNS -v ~/pkg:/go/pkg --rm centos7
	docker-compose run -v `pwd`:/go/src/github.com/STNS/STNS -v ~/pkg:/go/pkg --rm centos8
	docker-compose run -v `pwd`:/go/src/github.com/STNS/STNS -v ~/pkg:/go/pkg --rm ubuntu16
	docker-compose run -v `pwd`:/go/src/github.com/STNS/STNS -v ~/pkg:/go/pkg --rm ubuntu18
	docker-compose run -v `pwd`:/go/src/github.com/STNS/STNS -v ~/pkg:/go/pkg --rm debian8
	docker-compose run -v `pwd`:/go/src/github.com/STNS/STNS -v ~/pkg:/go/pkg --rm debian9

source_for_deb: ## Create source for DEB
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Distributing$(RESET)"
	rm -rf tmp.$(DIST) stns-v2-$(VERSION).orig.tar.gz
	mkdir -p tmp.$(DIST)/stns-v2-$(VERSION)
	cp -r $(SOURCES) tmp.$(DIST)/stns-v2-$(VERSION)
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
