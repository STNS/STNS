TEST ?= $(shell go list ./... | grep -v vendor)
VERSION = $(shell cat version)
REVISION = $(shell git describe --always)

GO=GO111MODULE=on go
INFO_COLOR=\033[1;34m
RESET=\033[0m
BOLD=\033[1m

DIST ?= unknown
PREFIX=/usr
BINDIR=$(PREFIX)/sbin
SOURCES=Makefile go.mod go.sum version model api middleware stns stns.go stns.conf.sample rpm/stns_v2.initd rpm/stns_v2.logrotate rpm/stns_v2.systemd
BUILD=tmp/bin

default: build

ci: depsdev test lint integration ## Run test and more...

depsdev: ## Installing dependencies for development
	$(GO) get github.com/golang/lint/golint
	$(GO) get github.com/pierrre/gotestcover
	$(GO) get -u github.com/tcnksm/ghr
	$(GO) get github.com/mitchellh/gox

test: ## Run test
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Testing$(RESET)"
	$(GO) test -v $(TEST) -timeout=30s -parallel=4
	$(GO) test -race $(TEST)

lint: ## Exec golint
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Linting$(RESET)"
	golint -min_confidence 1.1 -set_exit_status $(TEST)

server: ## Run server
	$(GO) run github.com/STNS/STNS --logfile ./stns.log --pidfile ./stns.pid --config ./stns/test.toml server

ghr: ## Upload to Github releases without token check
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Releasing for Github$(RESET)"
	ghr -u stns v$(VERSION)-$(REVISION) pkg

dist: build ## Upload to Github releases
	@test -z $(GITHUB_TOKEN) || test -z $(GITHUB_API) || $(MAKE) ghr

integration: ## Run integration test after Server wakeup
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Integration Testing$(RESET)"
	./misc/server start
	$(GO) test $(VERBOSE) -integration $(TEST) $(TEST_OPTIONS)
	./misc/server stop

build: ## Build server
	$(GO) build -o $(BUILD)/stns

install: build ## Install
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Installing as Server$(RESET)"
	cp $(BUILD)/stns $(BINDIR)/stns

docker:
	docker build -t nss_develop .
	docker run --cap-add=SYS_PTRACE --security-opt="seccomp=unconfined" -v `pwd`:/go/src/github.com/STNS/STNS -w /go/src/github.com/STNS/STNS -it nss_develop /bin/bash

source_for_rpm: ## Create source for RPM
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Distributing$(RESET)"
	rm -rf tmp.$(DIST) stns_v2-$(VERSION).tar.gz
	mkdir -p tmp.$(DIST)/stns_v2-$(VERSION)
	cp -r $(SOURCES) tmp.$(DIST)/stns_v2-$(VERSION)
	cd tmp.$(DIST) && \
		tar cf stns_v2-$(VERSION).tar stns_v2-$(VERSION) && \
		gzip -9 stns_v2-$(VERSION).tar
	cp tmp.$(DIST)/stns_v2-$(VERSION).tar.gz ./builds
	rm -rf tmp.$(DIST)

rpm: source_for_rpm ## Packaging for RPM
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Packaging for RPM$(RESET)"
	cp builds/stns_v2-$(VERSION).tar.gz /root/rpmbuild/SOURCES
	spectool -g -R rpm/stns.spec
	rpmbuild -ba rpm/stns.spec
	cp /root/rpmbuild/RPMS/*/*.rpm /go/src/github.com/STNS/STNS/builds


.PHONY: default dist test docker rpm source_for_rpm
