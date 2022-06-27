GO_BIN := /usr/local/bin/go

CHAIN_ID := testing
CHAIN_DIR := /tmp/trash
CHAIN_BIN := ./../../build/persistenceCore

all: clean setup start

.bash:
	CHAIN_ID=$(CHAIN_ID) \
	HOME=$(CHAIN_DIR) \
	CHAIN_BIN=$(CHAIN_BIN) \
	bash $(SCRIPT_FILE)

setup:
	$(MAKE) .bash SCRIPT_FILE=setup.sh

clean:
	rm -rf $(CHAIN_DIR)

start:
	HOME=$(CHAIN_DIR) $(CHAIN_BIN) start --minimum-gas-prices="0.0005stake"

###############################################################################
###                              Test commands                              ###
###############################################################################

run-contract:
	$(MAKE) .bash SCRIPT_FILE=contract.sh

run-upgrade:
	$(MAKE) .bash SCRIPT_FILE=upgrade.sh
