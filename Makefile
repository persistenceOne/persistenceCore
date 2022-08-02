#!/usr/bin/make -f

BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
COMMIT := $(shell git log -1 --format='%H')

# don't override user values
ifeq (,$(VERSION))
  VERSION := $(shell git describe --exact-match 2>/dev/null)
  # if VERSION is empty, then populate it with branch's name and raw commit hash
  ifeq (,$(VERSION))
    VERSION := $(BRANCH)-$(COMMIT)
  endif
endif

PACKAGES_SIMTEST=$(shell go list ./... | grep '/simulation')
LEDGER_ENABLED ?= true
SDK_PACK := $(shell go list -m github.com/cosmos/cosmos-sdk | sed  's/ /\@/g')
TM_VERSION := $(shell go list -m github.com/tendermint/tendermint | sed 's:.* ::') # grab everything after the space in "github.com/tendermint/tendermint v0.34.7"
BUILDDIR ?= $(CURDIR)/build

export GO111MODULE = on

include sims.mk

# process build tags
build_tags = netgo
ifeq ($(LEDGER_ENABLED),true)
  ifeq ($(OS),Windows_NT)
    GCCEXE = $(shell where gcc.exe 2> NUL)
    ifeq ($(GCCEXE),)
      $(error gcc.exe not installed for ledger support, please install or set LEDGER_ENABLED=false)
    else
      build_tags += ledger
    endif
  else
    UNAME_S = $(shell uname -s)
    ifeq ($(UNAME_S),OpenBSD)
      $(warning OpenBSD detected, disabling ledger support (https://github.com/cosmos/cosmos-sdk/issues/1988))
    else
      GCC = $(shell command -v gcc 2> /dev/null)
      ifeq ($(GCC),)
        $(error gcc not installed for ledger support, please install or set LEDGER_ENABLED=false)
      else
        build_tags += ledger
      endif
    endif
  endif
endif

build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

whitespace := $(subst ,, )
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

# process linker flags
ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=persistenceCore \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=persistenceCore \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X github.com/tendermint/tendermint/version.TMCoreSemVer=$(TM_VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep)

ifeq (cleveldb,$(findstring cleveldb,$(build_tags)))
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=cleveldb
endif
ifeq (linkstatic,$(findstring linkstatic,$(build_tags)))
  ldflags += -linkmode=external -extldflags "-Wl,-z,muldefs -static"
endif
ifeq (badgerdb,$(findstring badgerdb,$(build_tags)))
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=badgerdb
endif
ifeq (rocksdb,$(findstring rocksdb,$(build_tags)))
  CGO_ENABLED=1
endif

BUILD_FLAGS += -ldflags '${ldflags}' -tags "${build_tags}"

GOBIN = $(shell go env GOPATH)/bin
GOARCH = $(shell go env GOARCH)
GOOS = $(shell go env GOOS)

# Docker variables
DOCKER := $(shell which docker)

DOCKER_IMAGE_NAME = persistenceone/persistencecore
DOCKER_TAG_NAME = latest
DOCKER_CONTAINER_NAME = persistence-core-container
DOCKER_CMD ?= "/bin/sh"
DOCKER_VOLUME = -v $(CURDIR):/usr/local/app
DOCKER_FILE ?= docker/Dockerfile

.PHONY: all install build verify docker-run

###############################################################################
###                              Documentation                              ###
###############################################################################

all: install

BUILD_TARGETS := build install

build: BUILD_ARGS=-o $(BUILDDIR)/

$(BUILD_TARGETS): go.sum $(BUILDDIR)/
	go $@ -mod=readonly $(BUILD_FLAGS) $(BUILD_ARGS) ./...

$(BUILDDIR)/:
	mkdir -p $(BUILDDIR)/

build-linux: go.sum
	LEDGER_ENABLED=false GOOS=linux GOARCH=amd64 $(MAKE) build

go-mod-cache: go.sum
	@echo "--> Download go modules to local cache"
	@go mod download

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify

draw-deps:
	@# requires brew install graphviz or apt-get install graphviz
	go get github.com/RobotsAndPencils/goviz
	@goviz -i ./cmd/persistenceCore -d 2 | dot -Tpng -o dependency-graph.png

clean:
	rm -rf $(BUILDDIR)/ artifacts/

distclean: clean
	rm -rf vendor/

###############################################################################
###                              Proto                              		###
###############################################################################

proto-gen:
	@echo "Generating Protobuf files"
	scripts/protocgen.sh

###############################################################################
###                              Docker                             		###
###############################################################################

# Commands for running docker
#
# Run persistenceCore on docker
# Example Usage:
# 	make docker-build   ## Builds persistenceCore binary in 2 stages, 1st builder 2nd Runner
# 						   Final image only has the compiled persistenceCore binary
# 	make docker-interactive   ## Will start an shell session into the docker container
# 								 Access to persistenceCore binary here
# 		NOTE: To be used for testing only, since the container will be removed after stopping
# 	make docker-run DOCKER_CMD=sleep 10000000 DOCKER_OPTS=-d   ## Will run the container in the background
# 		NOTE: Recommeded to use docker commands directly for long running processes
# 	make docker-clean  # Will clean up the running container, as well as delete the image
# 						 after one is done testing
docker-build:
	$(DOCKER) buildx build ${DOCKER_ARGS} \
		-f $(DOCKER_FILE) \
		-t ${DOCKER_IMAGE_NAME}:${DOCKER_TAG_NAME} .

docker-build-push: docker-build
	$(DOCKER) push ${DOCKER_IMAGE_NAME}:${DOCKER_TAG_NAME}

docker-run:
	$(DOCKER) run --rm ${DOCKER_OPTS} ${DOCKER_VOLUME} --name=${DOCKER_CONTAINER_NAME} ${DOCKER_IMAGE_NAME}:${DOCKER_TAG_NAME} ${DOCKER_CMD}

docker-interactive:
	$(MAKE) docker-run DOCKER_CMD=/bin/bash DOCKER_OPTS="-it"

docker-clean-container:
	-$(DOCKER) stop ${DOCKER_CONTAINER_NAME}
	-$(DOCKER) rm ${DOCKER_CONTAINER_NAME}

docker-clean-image:
	-$(DOCKER) rmi ${DOCKER_IMAGE_NAME}:${DOCKER_TAG_NAME}

docker-clean: docker-clean-container docker-clean-image


###############################################################################
###                            Release commands                             ###
###############################################################################

PLATFORM ?= amd64

release-build-platform:
	@mkdir -p release/
	-@$(DOCKER) rm -f release-$(PLATFORM)
	$(MAKE) docker-build DOCKER_FILE="docker/Dockerfile.release" DOCKER_ARGS="--platform linux/$(PLATFORM) --no-cache" DOCKER_TAG_NAME="release-$(PLATFORM)"
	$(DOCKER) create -ti --name release-$(PLATFORM) ${DOCKER_IMAGE_NAME}:release-$(PLATFORM)
	$(DOCKER) cp release-$(PLATFORM):/usr/local/app/build/persistenceCore release/persistenceCore-$(VERSION)-linux-$(PLATFORM)
	tar -zcvf release/persistenceCore-$(VERSION)-linux-$(PLATFORM).tar.gz release/persistenceCore-$(VERSION)-linux-$(PLATFORM)
	-@$(DOCKER) rm -f release-$(PLATFORM)

release-sha:
	mkdir -p release/
	rm -f release/sha256sum.txt
	sha256sum release/* | sed 's#release/##g' > release/sha256sum.txt

# Create git archive
release-git:
	mkdir -p release/
	git archive \
		--format zip \
		--prefix "persistenceCore-$(VERSION)/" \
		-o "release/Source code.zip" \
		HEAD

	git archive \
		--format tar.gz \
		--prefix "persistenceCore-$(VERSION)/" \
		-o "release/Source code.tar.gz" \
		HEAD