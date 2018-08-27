TEST ?= $(shell go list ./... | grep -v vendor)
VERSION = $(shell cat version)
REVISION = $(shell git describe --always)

ifeq ("$(shell uname)","Darwin")
GO=GO111MODULE=on go
else
GO=GO111MODULE=on /usr/local/go/bin/go
endif
INFO_COLOR=\033[1;34m
RESET=\033[0m
BOLD=\033[1m

DIST ?= unknown
PREFIX=/usr
BINDIR=$(PREFIX)/sbin
SOURCES=Makefile go.mod go.sum version model api middleware stns stns.go package/

BUILD=tmp/bin

default: build

ci: depsdev test lint integration ## Run test and more...

depsdev: ## Installing dependencies for development
	$(GO) get github.com/golang/lint/golint
	$(GO) get -u github.com/tcnksm/ghr

test: ## Run test
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Testing$(RESET)"
	$(GO) test -v $(TEST) -timeout=30s -parallel=4
	$(GO) test -race $(TEST)

lint: ## Exec golint
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Linting$(RESET)"
	golint -min_confidence 1.1 -set_exit_status $(TEST)

server: ## Run server
	$(GO) run github.com/STNS/STNS --logfile ./stns.log --pidfile ./stns.pid --config ./stns/test.toml server

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
	docker run --cap-add=SYS_PTRACE --security-opt="seccomp=unconfined" -v $(GOPATH):/go/ -v $(GOPATH)/pkg/mod/cache:/go/pkg/mod/cache -w /go/src/github.com/STNS/STNS -it nss_develop /bin/bash

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
	docker-compose up

github_release: pkg ## Upload archives to Github Release on Mac
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Releasing for Github$(RESET)"
	rm -rf builds/.keep && ghr v$(VERSION) builds && git checkout builds/.keep

source_for_deb: ## Create source for DEB
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Distributing$(RESET)"
	rm -rf tmp.$(DIST) stns-v2-$(VERSION).orig.tar.gz
	mkdir -p tmp.$(DIST)/stns-v2-$(VERSION)
	cp -r $(SOURCES) tmp.$(DIST)/stns-v2-$(VERSION)
	cd tmp.$(DIST) && \
	tar zcf stns-v2-$(VERSION).tar.gz stns-v2-$(VERSION)
	mv tmp.$(DIST)/stns-v2-$(VERSION).tar.gz tmp.$(DIST)/stns-v2-$(VERSION).orig.tar.gz

deb_systemd: source_for_deb ## Packaging for DEB
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Packaging for DEB$(RESET)"
	cd tmp.$(DIST) && \
		tar xf stns-v2-$(VERSION).orig.tar.gz && \
		cd stns-v2-$(VERSION) && \
		dh_make --single --createorig -y && \
		rm -rf debian/*.ex debian/*.EX debian/README.Debian && \
		cp -r $(GOPATH)/src/github.com/STNS/STNS/debian-systemd/* debian/ && \
		sed -i -e 's/xenial/$(DIST)/g' debian/changelog && \
		debuild -uc -us
	cd tmp.$(DIST) && \
		find . -name "*.deb" | sed -e 's/\(\(.*stns-v2_.*\).deb\)/mv \1 \2.$(DIST).deb/g' | sh && \
		cp *.deb $(GOPATH)/src/github.com/STNS/STNS/builds
	rm -rf tmp.$(DIST)

.PHONY: default test docker rpm source_for_rpm pkg source_for_deb deb
