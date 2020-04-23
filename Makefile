PACKAGES := $(shell go list ./... | grep -v '/simulation')
VERSION := $(shell git branch | grep \* | cut -d ' ' -f2)
COMMIT := $(shell git rev-parse --short HEAD)
GOSUM := $(shell which gosum)

export GO111MODULE=on


BUILD_TAGS := -s  -w \
	-X github.com/persistenceOne/persistenceCore/version.Version=${VERSION} \
	-X github.com/persistenceOne/persistenceCore/version.Commit=${COMMIT}

ifneq (${GOSUM},)
	ifneq (${wildcard go.sum},)
		BUILD_TAGS += -X github.com/persistenceOne/persistenceCore/version.VendorHash=$(shell ${GOSUM} go.sum)
	endif
endif

BUILD_FLAGS += -ldflags "${BUILD_TAGS}"

all: install

build: go.sum
ifeq (${OS},Windows_NT)
	
	go build -mod=readonly ${BUILD_FLAGS} -o bin/coreClient.exe client/
	go build -mod=readonly ${BUILD_FLAGS} -o bin/coreNode.exe node/

else
	
	go build -mod=readonly ${BUILD_FLAGS} -o bin/coreClient client/
	go build -mod=readonly ${BUILD_FLAGS} -o bin/coreNode node/

endif

install: go.sum
	
	go install -mod=readonly ${BUILD_FLAGS} ./client
	go install -mod=readonly ${BUILD_FLAGS} ./node

go.sum:
	@echo "--> Ensuring dependencies have not been modified."
	@go mod verify

.PHONY: all build install  go.sum