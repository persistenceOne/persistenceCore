CHAIN_ID := testing

CHAIN_DIR ?= /tmp/trash
CHAIN_BIN ?= ./../../build/persistenceCore
CHAIN_DATA_DIR ?= .persistenceCore
WASM_PERMISSIONLESS ?= false

all: docker-clean clean setup start

.bash:
	CHAIN_ID=$(CHAIN_ID) \
	HOME=$(CHAIN_DIR) \
	CHAIN_BIN=$(CHAIN_BIN) \
	WASM_PERMISSIONLESS=$(WASM_PERMISSIONLESS) \
	CHAIN_DATA_DIR=$(CHAIN_DATA_DIR) \
	/bin/bash $(SCRIPT_FILE)

setup:
	$(MAKE) .bash SCRIPT_FILE=setup.sh

clean:
	rm -rf $(CHAIN_DIR)/$(CHAIN_DATA_DIR)

start:
	HOME=$(CHAIN_DIR) $(CHAIN_BIN) start --minimum-gas-prices="0.0005uxprt" $(ARGS)

###############################################################################
###                              Test commands                              ###
###############################################################################

run-wasm-contract:
	$(MAKE) .bash SCRIPT_FILE=wasm/contract.sh

run-wasm-gov-contract:
	$(MAKE) .bash SCRIPT_FILE=wasm/gov_contract.sh

run-wasm-cw20-base:
	$(MAKE) .bash SCRIPT_FILE=wasm/test_cw20.sh

run-wasm-cw20-govinit:
	$(MAKE) .bash SCRIPT_FILE=wasm/test_cw20_govinit.sh

run-upgrade:
	$(MAKE) .bash SCRIPT_FILE=upgrade.sh

###############################################################################
###                              Docker commands                            ###
###############################################################################

# Docker variables
DOCKER := $(shell which docker)

DOCKER_IMAGE_NAME = persistenceone/persistencecore
DOCKER_TAG_NAME = latest
DOCKER_CONTAINER := persistence-core-container
DOCKER_ENV ?= \
	-e HOME=/opt \
	-e CHAIN_DIR=/opt \
	-e CHAIN_BIN=/usr/bin/persistenceCore \
	-e WASM_PERMISSIONLESS=$(WASM_PERMISSIONLESS)
DOCKER_PORTS ?= \
	-p 127.0.0.1:1317:1317 \
	-p 127.0.0.1:26656-26657:26656-26657 \
	-p 127.0.0.1:9090:9090

docker-setup: docker-clean
	$(DOCKER) run --rm -d \
		--name=$(DOCKER_CONTAINER) \
		$(DOCKER_ENV) $(DOCKER_PORTS) \
		$(DOCKER_IMAGE_NAME):$(DOCKER_TAG_NAME) make
	echo "Waiting for the docker to start...."
	sleep 10

docker-exec:
	$(DOCKER) exec -it $(DOCKER_CONTAINER) /bin/bash

docker-clean:
	-$(DOCKER) stop $(DOCKER_CONTAINER)
