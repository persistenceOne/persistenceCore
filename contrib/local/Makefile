CHAIN_ID := testing
CHAIN_DIR := /tmp/trash
CHAIN_BIN := ./../../build/persistenceCore

setup:
	CHAIN_ID=$(CHAIN_ID) CHAIN_DIR=$(CHAIN_DIR) CHAIN_BIN=$(CHAIN_BIN) bash setup.sh

clean:
	rm -rf $(CHAIN_DIR)

start:
	HOME=$(CHAIN_DIR) $(CHAIN_BIN) start --minimum-gas-prices="0.0005stake"

###############################################################################
###                              Test commands                              ###
###############################################################################

software-upgrade:
	$(CHAIN_BIN) tx gov submit-proposal software-upgrade v2 --yes \
		--title "Test local upgrade" \
		--description "Test local upgrade" \
		--upgrade-height $(HEIGHT) \
		--from $(jq -r .genesis[0].address $(VALIDATOR_CONFIG)) \
	    --chain-id $(CHAIN_ID) \
		--node "https://rpc.testnet.persistence.one:443" \
		--deposit 600000000uxprt \
		--upgrade-info "The binary for this upgrade should be built from v0.2.x tagged or from master branch"