export GO111MODULE=on

VERSION := $(shell echo $(shell git describe --always) | sed 's/^v//')
COMMIT := $(shell git rev-parse --short HEAD)

build_tags = netgo
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=persistenceCore \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=persistenceCore \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep)

BUILD_FLAGS += -ldflags "${ldflags}"

GOBIN = $(shell go env GOPATH)/bin

all: verify build

install:
ifeq (${OS},Windows_NT)
	go build -mod=readonly ${BUILD_FLAGS} -o ${GOBIN}/persistenceCore.exe ./node

else
	go build -mod=readonly ${BUILD_FLAGS} -o ${GOBIN}/persistenceCore ./node

endif

build:
ifeq (${OS},Windows_NT)
	go build  ${BUILD_FLAGS} -o ${GOBIN}/persistenceCore.exe ./node

else
	go build  ${BUILD_FLAGS} -o ${GOBIN}/persistenceCore ./node

endif

verify:
	@echo "verifying modules"
	@go mod verify

.PHONY: all install build verify


DOCKER := $(shell which docker)

proto-gen:
	@echo "Generating Protobuf files"
	$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace persistence/sdk-proto-gen sh ./.script/protocgen.sh
