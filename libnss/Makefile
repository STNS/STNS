# The base of this code is https://github.com/linyows/stns/blob/master/Makefile
CC=gcc
CFLAGS=-Wall -Wstrict-prototypes -Werror -fPIC -std=c99 -D_GNU_SOURCE
LIBRARY=libnss_stns.so.2.0
LINKS=libnss_stns.so.2 libnss_stns.so
LD_SONAME=-Wl,-soname,libnss_stns.so.2

SRCS := $(wildcard *.c)
PREFIX=/usr
LIBDIR=$(PREFIX)/lib64
ifeq ($(wildcard $(LIBDIR)/.*),)
LIBDIR=$(PREFIX)/lib
endif
BINDIR=$(PREFIX)/bin

BUILD=tmp/libs
CRITERION_VERSION=2.3.2
SHUNIT_VERSION=2.1.6

ci: test integration_test
test: depsdev testdev ## Test with dependencies installation

build_dir: ## Create directory for build
	test -d $(BUILD) || mkdir -p $(BUILD)

depsdev: build_dir ## Installing dependencies for development
	test -f $(BUILD)/criterion.tar.bz2 || curl -sL https://github.com/Snaipe/Criterion/releases/download/v$(CRITERION_VERSION)/criterion-v$(CRITERION_VERSION)-linux-x86_64.tar.bz2 -o $(BUILD)/criterion.tar.bz2
	cd $(BUILD); tar xf criterion.tar.bz2; cd ../
	test -d /usr/include/criterion || mv $(BUILD)/criterion-v$(CRITERION_VERSION)/include/criterion /usr/include/criterion && mv $(BUILD)/criterion-v$(CRITERION_VERSION)/lib/libcriterion.* $(LIBDIR)/
	test -f $(BUILD)/shunit2.tgz || curl -sL https://storage.googleapis.com/google-code-archive-downloads/v2/code.google.com/shunit2/shunit2-$(SHUNIT_VERSION).tgz -o $(BUILD)/shunit2.tgz
	cd $(BUILD); tar xf shunit2.tgz; cd ../
	test -d /usr/include/shunit2 || mv $(BUILD)/shunit2-$(SHUNIT_VERSION)/ /usr/include/shunit2

testdev: ## Test without dependencies installation
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Testing$(RESET)"
	$(CC) -g3 -O0 $(SRCS) \
		-lcurl -ljansson -lcriterion -o $(BUILD)/test && \
		$(BUILD)/test --verbose
build: nss_build
nss_build: build_dir ## Build nss_stns
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Building nss_stns$(RESET)"
	$(CC) $(CFLAGS) -c toml.c -o $(BUILD)/toml.o
	$(CC) $(CFLAGS) -c stns_passwd.c -o $(BUILD)/stns_passwd.o
	$(CC) $(CFLAGS) -c stns_group.c -o $(BUILD)/stns_group.o
	$(CC) $(CFLAGS) -c stns_shadow.c -o $(BUILD)/stns_shadow.o
	$(CC) $(CFLAGS) -c stns.c -o $(BUILD)/stns.o
	$(CC) -shared $(LD_SONAME) -o $(BUILD)/$(LIBRARY) \
		$(BUILD)/stns.o \
		$(BUILD)/stns_passwd.o \
		$(BUILD)/toml.o \
		$(BUILD)/stns_group.o \
		$(BUILD)/stns_shadow.o \
		-lcurl -ljansson

integration_test: build install depsdev ## Run integration test
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Integration Testing$(RESET)"
	cd .. && misc/server start
	test -d /usr/lib/x86_64-linux-gnu && ln -sf /usr/lib/libnss_stns.so.2.0 /usr/lib/x86_64-linux-gnu/libnss_stns.so.2.0 || true
	mkdir -p /etc/stns/client
	cp stns.conf.example /etc/stns/client/stns.conf
	sed -i 's/1104/8050/' /etc/stns/client/stns.conf
	sed -i -e 's/^passwd:.*/passwd: files stns/g' /etc/nsswitch.conf
	sed -i -e 's/^shadow:.*/shadow: files stns/g' /etc/nsswitch.conf
	sed -i -e 's/^group:.*/group: files stns/g' /etc/nsswitch.conf
	test/integration_test.sh
	cd .. && misc/server stop

install: install_lib ## Install stns

install_lib: ## Install only shared objects
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Installing as Libraries$(RESET)"
	[ -d $(LIBDIR) ] || install -d $(LIBDIR)
	install $(BUILD)/$(LIBRARY) $(LIBDIR)
	cd $(LIBDIR); for link in $(LINKS); do ln -sf $(LIBRARY) $$link ; done;

.PHONY: depsdev test testdev build
