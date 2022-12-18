FILE = config.yaml

HELM_REPO = shuttle
HELM_CHART = devnet
HELM_NAME = persistence
HELM_VERSION = 0.1.10

###############################################################################
###                              Helm commands                              ###
###############################################################################

setup:
	helm repo add $(HELM_REPO) https://anmol1696.github.io/shuttle/
	helm repo update
	helm search repo $(HELM_REPO)/$(HELM_CHART) --version $(HELM_VERSION)

install:
	helm install -f $(FILE) $(HELM_NAME) $(HELM_REPO)/$(HELM_CHART) --version $(HELM_VERSION)

upgrade:
	helm upgrade --debug -f $(FILE) $(HELM_NAME) $(HELM_REPO)/$(HELM_CHART) --version $(HELM_VERSION)

debug:
	helm install --dry-run --debug -f $(FILE) $(HELM_NAME) $(HELM_REPO)/$(HELM_CHART) --version $(HELM_VERSION)

delete:
	-helm delete $(HELM_NAME)

###############################################################################
###                              Test commands                              ###
###############################################################################
OFFSET_HEIGHT = 80
UPGRADE_NAME = v6
CHAIN_ID = $(shell yq -r ".chains[0].name" $(FILE))

test-upgrade:
	OFFSET_HEIGHT=$(OFFSET_HEIGHT) \
	UPGRADE_NAME=$(UPGRADE_NAME) \
	CHAIN_ID=$(CHAIN_ID) \
	bash -e upgrade.sh

setup-local-val:
	echo "Sets up local directory to a validator, which can be connected with the chain"
	echo "Not Implemented"

###############################################################################
###                              Port forward                              ###
###############################################################################

.PHONY: port-forward port-forward-all
.port-forward:
	kubectl port-forward pods/$(chain)-genesis-0 $(localrpc):26657 &
	kubectl port-forward pods/$(chain)-genesis-0 $(localp2p):26656 &
	kubectl port-forward pods/$(chain)-genesis-0 $(localrest):1317 &

num_chains = $(shell yq -r ".chains | length" $(FILE))
port-forward-all:
	echo "Port forwarding all chains to localhost"
	for i in $(shell seq 0 $(num_chains)); do \
  		$(MAKE) .port-forward \
  			chain=$$(yq -r ".chains[$$i].name" $(FILE)) \
  			localrpc=$$(yq -r ".chains[$$i].ports.rpc" $(FILE)) \
  			localrest=$$(yq -r ".chains[$$i].ports.rest" $(FILE)); \
	done
	echo "Port forwarding explorer to localhost"
	kubectl port-forward service/explorer 8080:8080 &

.PHONY: stop-forward
stop-forward:
	-pkill -f "port-forward"
