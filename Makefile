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
ifeq (badgerdb,$(findstring badgerdb,$(build_tags)))
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=badgerdb
endif
ifeq (rocksdb,$(findstring rocksdb,$(build_tags)))
  CGO_ENABLED=1
endif

BUILD_FLAGS += -ldflags "${ldflags}" -tags "${build_tags}"

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

.PHONY: all install build verify docker-run

###############################################################################
###                              Documentation                              ###
###############################################################################

all: install lint test

BUILD_TARGETS := build install

build: BUILD_ARGS=-o $(BUILDDIR)/

$(BUILD_TARGETS): go.sum $(BUILDDIR)/
	go $@ -mod=readonly $(BUILD_FLAGS) $(BUILD_ARGS) ./...

$(BUILDDIR)/:
	mkdir -p $(BUILDDIR)/

build-reproducible: go.sum
	$(DOCKER) rm latest-build || true
	$(DOCKER) run --volume=$(CURDIR):/sources:ro \
        --env TARGET_PLATFORMS='linux/amd64 darwin/amd64 linux/arm64 windows/amd64' \
        --env APP=persistenceCore \
        --env VERSION=$(VERSION) \
        --env COMMIT=$(COMMIT) \
        --env LEDGER_ENABLED=$(LEDGER_ENABLED) \
        --name latest-build tendermintdev/rbuilder:latest
	$(DOCKER) cp -a latest-build:/home/builder/artifacts/ $(CURDIR)/

build-linux: go.sum
	LEDGER_ENABLED=false GOOS=linux GOARCH=amd64 $(MAKE) build

build-contract-tests-hooks:
	mkdir -p $(BUILDDIR)
	go build -mod=readonly $(BUILD_FLAGS) -o $(BUILDDIR)/ ./cmd/contract_tests

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

ifeq (${OS},Windows_NT)
	bin_name = persistenceCore
else
	bin_name = persistenceCore.exe
endif

release: build
	mkdir -p release
	tar -czvf release/persistenceCore-${GOOS}-${GOARCH}.tar.gz --directory=build/${GOOS}/${GOARCH} ${bin_name}

###############################################################################
###                              Proto                              		###
###############################################################################

proto-gen:
	@echo "Generating Protobuf files"
	$(DOCKER) run --rm \
		-v $(shell go list -f "{{ .Dir }}" -m github.com/cosmos/cosmos-sdk):/workspace/cosmos_sdk_dir \
	 	--env COSMOS_SDK_DIR=/workspace/cosmos_sdk_dir \
	 	-v $(CURDIR):/workspace --workdir /workspace \
		tendermintdev/sdk-proto-gen sh ./scripts/protocgen.sh

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
	${DOCKER} build -t ${DOCKER_IMAGE_NAME}:${DOCKER_TAG_NAME} .

docker-build-push: docker-build
	${DOCKER} push ${DOCKER_IMAGE_NAME}:${DOCKER_TAG_NAME}

docker-run:
	${DOCKER} run ${DOCKER_OPTS} ${DOCKER_VOLUME} --name=${DOCKER_CONTAINER_NAME} ${DOCKER_IMAGE_NAME}:${DOCKER_TAG_NAME} ${DOCKER_CMD}

docker-interactive:
	${MAKE} docker-run DOCKER_CMD=/bin/sh DOCKER_OPTS="--rm -it"

docker-clean-container:
	-${DOCKER} stop ${DOCKER_CONTAINER_NAME}
	-${DOCKER} rm ${DOCKER_CONTAINER_NAME}

docker-clean-image:
	-${DOCKER} rmi ${DOCKER_IMAGE_NAME}:${DOCKER_TAG_NAME}

docker-clean: docker-clean-container docker-clean-image
