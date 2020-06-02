export GO111MODULE=on

VERSION := $(shell git branch | grep \* | cut -d ' ' -f2)
COMMIT := $(shell git rev-parse --short HEAD)

BUILD_TAGS := -s  -w \
	-X github.com/persistenceOne/assetMantle/version.Version=${VERSION} \
	-X github.com/persistenceOne/assetMantle/version.Commit=${COMMIT}

BUILD_FLAGS += -ldflags "${BUILD_TAGS}"

all: verify build

install:
ifeq (${OS},Windows_NT)
	
	go build -mod=readonly ${BUILD_FLAGS} -o ${GOBIN}/assetClient.exe ./client
	go build -mod=readonly ${BUILD_FLAGS} -o ${GOBIN}/assetNode.exe ./node

else
	
	go build -mod=readonly ${BUILD_FLAGS} -o ${GOBIN}/assetClient ./client
	go build -mod=readonly ${BUILD_FLAGS} -o ${GOBIN}/assetNode ./node

endif

build:
ifeq (${OS},Windows_NT)

	go build  ${BUILD_FLAGS} -o ${GOBIN}/assetClient.exe ./client
	go build  ${BUILD_FLAGS} -o ${GOBIN}/assetNode.exe ./node

else

	go build  ${BUILD_FLAGS} -o ${GOBIN}/assetClient ./client
	go build  ${BUILD_FLAGS} -o ${GOBIN}/assetNode ./node

endif

verify:
	@echo "verifying modules"
	@go mod verify

.PHONY: all install build verify