PACKAGES := $(shell go list ./... | grep -v '/simulation')
VERSION := $(shell git branch | grep \* | cut -d ' ' -f2)
COMMIT := $(shell git rev-parse --short HEAD)
GOSUM := $(shell which gosum)

export GO111MODULE=on


BUILD_TAGS := -s  -w \
	-X github.com/persistenceOne/persistenceSDK/version.Version=${VERSION} \
	-X github.com/persistenceOne/persistenceSDK/version.Commit=${COMMIT}

ifneq (${GOSUM},)
	ifneq (${wildcard go.sum},)
		BUILD_TAGS += -X github.com/persistenceOne/persistenceSDK/version.VendorHash=$(shell ${GOSUM} go.sum)
	endif
endif

BUILD_FLAGS += -ldflags "${BUILD_TAGS}"

all: install

build: go.sum
ifeq (${OS},Windows_NT)
	
	go build -mod=readonly ${BUILD_FLAGS} -o bin/hubClient.exe commands/hub/hubClient/
	go build -mod=readonly ${BUILD_FLAGS} -o bin/hubNode.exe commands/hub/hubNode/
	
	go build -mod=readonly ${BUILD_FLAGS} -o bin/zoneClient.exe commands/zone/zoneClient/
	go build -mod=readonly ${BUILD_FLAGS} -o bin/zoneNode.exe commands/zone/zoneNode/

else
	
	go build -mod=readonly ${BUILD_FLAGS} -o bin/hubClient commands/hub/hubClient/
	go build -mod=readonly ${BUILD_FLAGS} -o bin/hubNode commands/hub/hubNode/
	
	go build -mod=readonly ${BUILD_FLAGS} -o bin/zoneClient commands/zone/zoneClient/
	go build -mod=readonly ${BUILD_FLAGS} -o bin/zoneNode commands/zone/zoneNode/

endif

install: go.sum
	
	go install -mod=readonly ${BUILD_FLAGS} ./commands/hub/hubClient
	go install -mod=readonly ${BUILD_FLAGS} ./commands/hub/hubNode

	go install -mod=readonly ${BUILD_FLAGS} ./commands/zone/zoneClient
	go install -mod=readonly ${BUILD_FLAGS} ./commands/zone/zoneNode

go.sum:
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify

.PHONY: all build install  go.sum