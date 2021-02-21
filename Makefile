.PHONY: install

TOPDIR = $(shell pwd)

export GOWORKSPACE := $(shell pwd)
export GOBIN := $(GOWORKSPACE)/bin
export GO111MODULE := on

GIT_VERSION = $(shell bash debian/print-git-version.sh)

GO_BINS =
GO_BINS += git-caddy

all:
	$(MAKE) compile

compile:
	mkdir -p tmp
	mkdir -p $(GOBIN)
	go install github.com/...

install:
	mkdir -p $(DESTDIR)/$(INSTALL_PREFIX)/bin/ 
	mkdir -p $(DESTDIR)/$(INSTALL_PREFIX)/data/
	mkdir -p $(DESTDIR)/$(INSTALL_PREFIX)/scripts/
	mkdir -p $(DESTDIR)/opt/polaris-agent/plugins
	cp -v $(GOWORKSPACE)/bin-extra/* $(DESTDIR)/$(INSTALL_PREFIX)/bin/
	cp -v scripts/pagcli.py $(DESTDIR)/$(INSTALL_PREFIX)/scripts/.py
	cp -vr data/* $(DESTDIR)/$(INSTALL_PREFIX)/data/
	$(MAKE) install-bins

install-bins: $(addprefix installbin-, $(GO_BINS))

$(addprefix installbin-, $(GO_BINS)):
	$(eval BIN=$(@:installbin-%=%))
	cp -v $(GOBIN)/$(BIN) $(DESTDIR)/$(INSTALL_PREFIX)/bin/

clean:
	$(MAKE) -C bindata clean
