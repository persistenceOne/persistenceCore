UPGRADE_NAME = v7
OFFSET_HEIGHT = 50
# index of the chain to be upgraded
UPGRADE_CHAIN_INDEX = 0

FILE = $(UPGRADE_NAME)/config.yaml

HELM_REPO = starship
HELM_CHART = devnet
HELM_NAME = persistence
HELM_VERSION = 0.1.15
# k8s namespace
NAMESPACE = dev-core

# max chain index
num_chains = $(shell yq -r ".chains | length - 1" $(FILE))

###############################################################################
###                              Helm commands                              ###
###############################################################################

setup-helm:
	helm repo add $(HELM_REPO) https://anmol1696.github.io/starship/
	helm repo update
	helm search repo $(HELM_REPO)/$(HELM_CHART) --version $(HELM_VERSION)

install-helm: delete-helm
	helm install -f $(FILE) $(HELM_NAME) $(HELM_REPO)/$(HELM_CHART) --version $(HELM_VERSION) --namespace $(NAMESPACE)

upgrade-helm:
	helm upgrade --debug -f $(FILE) $(HELM_NAME) $(HELM_REPO)/$(HELM_CHART) --version $(HELM_VERSION) --namespace $(NAMESPACE)

debug-helm:
	helm install --dry-run --debug -f $(FILE) $(HELM_NAME) $(HELM_REPO)/$(HELM_CHART) --version $(HELM_VERSION) --namespace $(NAMESPACE)

delete-helm:
	-helm delete $(HELM_NAME) --namespace $(NAMESPACE)

###############################################################################
###                          Test upgrade commands                          ###
###############################################################################

define chain-env
	CHAIN_ID=$(shell yq -r ".chains[$(1)].name" $(FILE)) \
	CHAIN_BIN=$(shell yq -r ".chains[$(1)].bin" $(FILE)) \
	CHAIN_DATA_DIR=$(shell yq -r ".chains[$(1)].dir" $(FILE)) \
	NODE_PORT=$(shell yq -r ".chains[$(1)].ports.rpc" $(FILE)) \
	EXPOSER_PORT=$(shell yq -r ".chains[$(1)].ports.exposer" $(FILE))
endef

.setup-keys:
	@$(call chain-env, $(index)) bash -e ./setup.sh

setup-keys:
	for i in {0..$(num_chains)}; do \
		$(MAKE) .setup-keys index=$$i; \
	done

.clean-keys:
	@$(call chain-env, $(index)) && \
	echo "deleting ~/$$CHAIN_DATA_DIR/keyring-test" && \
	rm -rf ~/$$CHAIN_DATA_DIR/keyring-test

clean-keys:
	for i in {0..$(num_chains)}; do \
		$(MAKE) .clean-keys index=$$i; \
	done

.reset-data-dir:
	@$(call chain-env, $(index)) && \
	echo "deleting & recreating ~/$$CHAIN_DATA_DIR" && \
	rm -rf ~/$$CHAIN_DATA_DIR && \
	mkdir ~/$$CHAIN_DATA_DIR

reset-data-dir:
	for i in {0..$(num_chains)}; do \
		$(MAKE) .reset-data-dir index=$$i; \
	done

pre-upgrade:
	bash -e $(UPGRADE_NAME)/pre-upgrade.sh

upgrade-chain:
	@$(call chain-env, $(UPGRADE_CHAIN_INDEX)) \
	UPGRADE_NAME=$(UPGRADE_NAME) \
	OFFSET_HEIGHT=$(OFFSET_HEIGHT) \
	bash -e ./upgrade.sh

post-upgrade:
	bash -e $(UPGRADE_NAME)/post-upgrade.sh

###############################################################################
###                              Port forward                               ###
###############################################################################

.PHONY: port-forward port-forward-all
.port-forward:
	kubectl port-forward pods/$(chain)-genesis-0 $(localrpc):26657 --namespace $(NAMESPACE) &
	kubectl port-forward pods/$(chain)-genesis-0 $(localp2p):26656 --namespace $(NAMESPACE) &
	kubectl port-forward pods/$(chain)-genesis-0 $(localrest):1317 --namespace $(NAMESPACE) &
	kubectl port-forward pods/$(chain)-genesis-0 $(localexposer):8081 --namespace $(NAMESPACE) &

port-forward-exposer:
	for i in $(shell seq 0 $(num_vals)); do \
		kubectl port-forward pods/$(chain)-validator-$$(i) `expr $(localexposer) + $i + 1`:8080 --namespace $(NAMESPACE) &; \
	done

port-forward-all:
	echo "Port forwarding all chains to localhost"
	for i in $(shell seq 0 $(num_chains)); do \
  		$(MAKE) .port-forward \
  			chain=$$(yq -r ".chains[$$i].name" $(FILE)) \
  			localrpc=$$(yq -r ".chains[$$i].ports.rpc" $(FILE)) \
  			localexposer=$$(yq -r ".chains[$$i].ports.exposer" $(FILE)) \
  			localrest=$$(yq -r ".chains[$$i].ports.rest" $(FILE)); \
	done
	echo "Port forwarding explorer to localhost"
	kubectl port-forward service/explorer 8080:8080 --namespace $(NAMESPACE) &

port-forward-all-exposers:
	echo "Port forwarding all chains exposers to localhost"
	for i in $(shell seq 0 $(num_chains)); do \
		$(MAKE) port-forward-exposer \
			chain=$$(yq -r ".chains[$$i].name" $(FILE)) \
			localexposer=$$(yq -r ".chains[$$i].ports.exposer" $(FILE)) \
			num_vals=$$(yq -r ".chains[$$i].numValidators" $(FILE)); \
	done

.PHONY: stop-forward
stop-forward:
	-@pkill -f "port-forward"
