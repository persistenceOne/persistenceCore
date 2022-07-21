CHAIN_ID := testing

CHAIN_DIR ?= /tmp/trash
CHAIN_BIN ?= ./../../build/persistenceCore
WASM_PERMISSIONLESS ?= false

all: clean setup start

.bash:
	CHAIN_ID=$(CHAIN_ID) \
	HOME=$(CHAIN_DIR) \
	CHAIN_BIN=$(CHAIN_BIN) \
	WASM_PERMISSIONLESS=$(WASM_PERMISSIONLESS) \
	bash $(SCRIPT_FILE)

setup:
	$(MAKE) .bash SCRIPT_FILE=setup.sh

clean:
	rm -rf $(CHAIN_DIR)/.persistenceCore

start:
	HOME=$(CHAIN_DIR) $(CHAIN_BIN) start --minimum-gas-prices="0.0005stake" $(ARGS)

###############################################################################
###                              Test commands                              ###
###############################################################################

run-contract:
	$(MAKE) .bash SCRIPT_FILE=contract.sh

run-gov-contract:
	$(MAKE) .bash SCRIPT_FILE=gov_contract.sh

run-cw20-test:
	$(MAKE) .bash SCRIPT_FILE=test_cw20.sh

run-test:
	$(MAKE) .bash SCRIPT_FILE=test.sh

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

docker-setup: docker-clean
	$(DOCKER) run --rm -d \
		--name=$(DOCKER_CONTAINER) \
		$(DOCKER_ENV) \
		$(DOCKER_IMAGE_NAME):$(DOCKER_TAG_NAME) make

docker-exec:
	$(DOCKER) exec -it $(DOCKER_CONTAINER) /bin/bash

docker-clean:
	-$(DOCKER) stop $(DOCKER_CONTAINER)
