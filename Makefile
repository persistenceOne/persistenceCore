export GO111MODULE=on

VERSION := $(shell git branch | grep \* | cut -d ' ' -f2)
COMMIT := $(shell git rev-parse --short HEAD)

BUILD_TAGS := -s  -w \
	-X github.com/persistenceOne/persistenceCore/version.Version=${VERSION} \
	-X github.com/persistenceOne/persistenceCore/version.Commit=${COMMIT}

BUILD_FLAGS += -ldflags "${BUILD_TAGS}"

all: verify build

install:
ifeq (${OS},Windows_NT)
	
	go build -mod=readonly ${BUILD_FLAGS} -o ${GOBIN}/coreClient.exe ./client
	go build -mod=readonly ${BUILD_FLAGS} -o ${GOBIN}/coreNode.exe ./node

else
	
	go build -mod=readonly ${BUILD_FLAGS} -o ${GOBIN}/coreClient ./client
	go build -mod=readonly ${BUILD_FLAGS} -o ${GOBIN}/coreNode ./node

endif

build:
ifeq (${OS},Windows_NT)

	go build  ${BUILD_FLAGS} -o ${GOBIN}/coreClient.exe ./client
	go build  ${BUILD_FLAGS} -o ${GOBIN}/coreNode.exe ./node

else

	go build  ${BUILD_FLAGS} -o ${GOBIN}/coreClient ./client
	go build  ${BUILD_FLAGS} -o ${GOBIN}/coreNode ./node

endif

verify:
	@echo "verifying modules"
	@go mod verify

.PHONY: all install build verify