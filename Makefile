#!/usr/bin/make -f

BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
COMMIT := $(shell git log -1 --format='%H')

# don't override user values
ifeq (,$(VERSION))
  VERSION := $(shell git describe --tags --exact-match)
  # if VERSION is empty, then populate it with branch's name and raw commit hash
  ifeq (,$(VERSION))
    VERSION := $(BRANCH)-$(COMMIT)
  endif
endif

BUILDDIR ?= $(CURDIR)/build
PACKAGES_SIMTEST=$(shell go list ./... | grep '/simulation')
LEDGER_ENABLED ?= true
SDK_PACK := $(shell go list -m github.com/cosmos/cosmos-sdk | sed  's/ /\@/g')
BFT_VERSION := $(shell go list -m github.com/cometbft/cometbft | sed 's:.* ::') # grab everything after the space in "github.com/cometbft/cometbft v0.34.7"

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

whitespace :=
whitespace += $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

# process linker flags

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=persistenceCore \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=persistenceCore \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep)" \
		  -X github.com/cometbft/cometbft/version.TMCoreSemVer=$(BFT_VERSION)

# DB backend selection
ifeq (cleveldb,$(findstring cleveldb,$(CORE_BUILD_OPTIONS)))
  build_tags += gcc
endif
ifeq (badgerdb,$(findstring badgerdb,$(CORE_BUILD_OPTIONS)))
  build_tags += badgerdb
endif
# handle rocksdb
ifeq (rocksdb,$(findstring rocksdb,$(CORE_BUILD_OPTIONS)))
  CGO_ENABLED=1
  build_tags += rocksdb
endif
# handle boltdb
ifeq (boltdb,$(findstring boltdb,$(CORE_BUILD_OPTIONS)))
  build_tags += boltdb
endif

ifeq (,$(findstring nostrip,$(CORE_BUILD_OPTIONS)))
  ldflags += -w -s
endif
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'
# check for nostrip option
ifeq (,$(findstring nostrip,$(CORE_BUILD_OPTIONS)))
  BUILD_FLAGS += -trimpath
endif

# Check for debug option
ifeq (debug,$(findstring debug,$(CORE_BUILD_OPTIONS)))
  BUILD_FLAGS += -gcflags "all=-N -l"
endif

# Docker variables
DOCKER := $(shell which docker)

include sims.mk

###############################################################################
###                                  Build                                  ###
###############################################################################

all: build lint

install: go.sum
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/persistenceCore

build:
	go build $(BUILD_FLAGS) -o bin/persistenceCore ./cmd/persistenceCore

$(BUILDDIR)/:
	mkdir -p $(BUILDDIR)/

vulncheck: $(BUILDDIR)/
	GOBIN=$(BUILDDIR) go install golang.org/x/vuln/cmd/govulncheck@latest
	$(BUILDDIR)/govulncheck ./...

.PHONY: all install lint build vulncheck

###############################################################################
###                          Tools & Dependencies                           ###
###############################################################################

go.sum: go.mod
	@echo "Ensure dependencies have not been modified ..." >&2
	go mod verify
	go mod tidy

###############################################################################
###                                Linting                                  ###
###############################################################################

golangci_lint_cmd=golangci-lint
golangci_version=v1.53.3

lint:
	@echo "--> Running linter"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(golangci_version)
	@$(golangci_lint_cmd) run --timeout=10m

lint-fix:
	@echo "--> Running linter"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(golangci_version)
	@$(golangci_lint_cmd) run --fix --out-format=tab --issues-exit-code=0

.PHONY: lint lint-fix

###############################################################################
###                              Documentation                              ###
###############################################################################

draw-deps:
	@# requires brew install graphviz or apt-get install graphviz
	go get github.com/RobotsAndPencils/goviz
	@goviz -i ./cmd/persistenceCore -d 2 | dot -Tpng -o dependency-graph.png


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
# 		NOTE: Recommended to use docker commands directly for long running processes
# 	make docker-clean  # Will clean up the running container, as well as delete the image
# 						 after one is done testing

include docker/Makefile


###############################################################################
###                            Release commands                             ###
###############################################################################

PLATFORM ?= amd64

release-build-platform:
	@mkdir -p release/
	-@$(DOCKER) rm -f release-$(PLATFORM)
	$(MAKE) docker-build PROCESS="persistencecore" DOCKER_FILE="Dockerfile.release" \
		DOCKER_BUILD_ARGS="--platform linux/$(PLATFORM) --no-cache --load" \
		DOCKER_TAG_NAME="release-$(PLATFORM)"
	$(DOCKER) images
	$(DOCKER) create -ti --name release-$(PLATFORM) $(DOCKER_IMAGE_NAME):release-$(PLATFORM)
	$(DOCKER) cp release-$(PLATFORM):/usr/local/app/bin/persistenceCore release/persistenceCore-$(VERSION)-linux-$(PLATFORM)
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


###############################################################################
###                   Docker Build (heighliner)                             ###
###############################################################################

get-heighliner:
	git clone https://github.com/strangelove-ventures/heighliner.git
	cd heighliner && go install

local-image:
ifeq (,$(shell which heighliner))
	echo 'heighliner' binary not found. Consider running `make get-heighliner`
else
	heighliner build -c persistence --local -f ./chains.yaml
endif

.PHONY: get-heighliner local-image

###############################################################################
###                                   testing                               ###
###############################################################################
test: ictest-all

# TODO: add runsim and benchmarking

###############################################################################
###                             e2e interchain test                         ###
###############################################################################

ictest-all: rm-testcache
	cd interchaintest && go test -v -run ./...

# Executes basic chain test via interchaintest
ictest-basic: rm-testcache
	cd interchaintest && go test -race -v -run TestBasicPersistenceStart .

ictest-ibchooks: rm-testcache
	cd interchaintest && go test -race -v -run TestPersistenceIBCHooks .

ictest-pfm: rm-testcache
	cd interchaintest && go test -race -v -run TestPacketForwardMiddlewareRouter .

# Executes a chain upgrade test via interchaintest
ictest-upgrade: rm-testcache
	cd interchaintest && go test -race -v -run TestPersistenceUpgradeBasic .

# Executes a chain upgrade locally via interchaintest after compiling a local image as persistence:local
ictest-upgrade-local: local-image ictest-upgrade

# Executes IBC tests via interchaintest
ictest-ibc: rm-testcache
	cd interchaintest && go test -race -v -run TestPersistenceGaiaIBCTransfer .

# Executes Skip's MEV auction module tests via interchaintest
ictest-pob: rm-testcache
	cd interchaintest && go test -race -v -run TestSkipMevAuction .

# Executes LSM tests
ictest-lsm: rm-testcache
	cd interchaintest && go test -race -v -run "(TestMultiTokenizeVote|TestTokenizeSendVote|TestBondTokenize)" .

ictest-haltfork: rm-testcache
	cd interchaintest && go test -race -v -run TestPersistenceLSMHaltFork .

# Executes the main Liquidstake test
ictest-liquidstake: rm-testcache
	cd interchaintest && go test -race -v -run TestLiquidStakeStkXPRT .

# Executes all Liquidstake tests
ictest-liquidstake-all: rm-testcache
	cd interchaintest && go test -race -v -run "(TestLiquidStakeStkXPRT|TestLiquidStakeUnstakeStkXPRT|TestPauseLiquidStakeStkXPRT)" .

rm-testcache:
	go clean -testcache

.PHONY: test ictest-all ictest-basic ictest-ibchooks ictest-pfm ictest-upgrade ictest-upgrade-local ictest-ibc ictest-pob ictest-lsm ictest-liquidstake
