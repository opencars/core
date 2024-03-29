.PHONY: default build clean
APPS        := grpc-server
BLDDIR      ?= bin
VERSION     ?= $(shell cat VERSION)
IMPORT_BASE := github.com/opencars/core
LDFLAGS     := -ldflags "-X $(IMPORT_BASE)/pkg/version.Version=$(VERSION)"

default: clean build

build: $(APPS)

$(BLDDIR)/%:
	go build --race $(LDFLAGS) -o $@ ./cmd/$*

$(APPS): %: $(BLDDIR)/%

clean:
	@mkdir -p $(BLDDIR)
	@for app in $(APPS) ; do \
		rm -f $(BLDDIR)/$$app ; \
	done
