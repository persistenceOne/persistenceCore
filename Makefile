export GO111MODULE=on

VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
TM_VERSION := $(shell go list -m github.com/tendermint/tendermint | sed 's:.* ::')
COMMIT := $(shell git rev-parse --short HEAD)
LEDGER_ENABLED ?= true
include sims.mk

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

all: verify build

install:
ifeq (${OS},Windows_NT)
	go build -mod=readonly ${BUILD_FLAGS} -o ${GOBIN}/persistenceCore.exe ./node

else
	go build -mod=readonly ${BUILD_FLAGS} -o ${GOBIN}/persistenceCore ./node

endif

build:
ifeq (${OS},Windows_NT)
	go build  ${BUILD_FLAGS} -o build/${GOOS}/${GOARCH}/persistenceCore.exe ./node

else
	go build  ${BUILD_FLAGS} -o build/${GOOS}/${GOARCH}/persistenceCore ./node

endif

verify:
	@echo "verifying modules"
	@go mod verify

release: build
	mkdir -p release
ifeq (${OS},Windows_NT)
	tar -czvf release/persistenceCore-${GOOS}-${GOARCH}.tar.gz --directory=build/${GOOS}/${GOARCH} persistenceCore.exe
else
	tar -czvf release/persistenceCore-${GOOS}-${GOARCH}.tar.gz --directory=build/${GOOS}/${GOARCH} persistenceCore
endif


clean:
	rm -rf build release

proto-gen:
	@echo "Generating Protobuf files"
	$(DOCKER) run --rm -v $(shell go list -f "{{ .Dir }}" \
	-m github.com/cosmos/cosmos-sdk):/workspace/cosmos_sdk_dir\
	 --env COSMOS_SDK_DIR=/workspace/cosmos_sdk_dir \
	 -v $(CURDIR):/workspace --workdir /workspace \
	 tendermintdev/sdk-proto-gen sh ./.script/protocgen.sh


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
