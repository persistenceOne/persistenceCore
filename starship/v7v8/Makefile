HELM_NAME = update-test
HELM_FILE = config.yaml

HELM_REPO = starship
HELM_CHART = devnet
HELM_VERSION = 0.1.37

###############################################################################
###                              Staringup                                  ###
###############################################################################

.PHONY: start
start: setup-helm install port-forward

.PHONY: stop
stop: stop-forward delete

###############################################################################
###                              Helm Charts                                ###
###############################################################################

setup-helm:
	helm repo add $(HELM_REPO) https://raw.githubusercontent.com/persistenceOne/starship/ajeet/sdk-v47/helm_repo
	helm repo update
	helm search repo $(HELM_REPO)/$(HELM_CHART) --version $(HELM_VERSION)

install:
	helm install -f $(HELM_FILE) $(HELM_NAME) $(HELM_REPO)/$(HELM_CHART) --wait --debug --timeout=20m --version $(HELM_VERSION)

upgrade:
	helm upgrade --debug $(HELM_NAME) $(HELM_REPO)/$(HELM_CHART) -f $(HELM_FILE) --version $(HELM_VERSION)

debug:
	helm install --dry-run --debug -f $(HELM_FILE) $(HELM_NAME) $(HELM_REPO)/$(HELM_CHART)

delete:
	-helm delete $(HELM_NAME)

###############################################################################
###                              Tests cases                                ###
###############################################################################

.PHONY: test
test:
	go test . -short -v -count=1
	go test . -v -count=1
	go test . -short -v -count=1


###############################################################################
###                              Port forward                              ###
###############################################################################

.PHOY: port-forward
port-forward:
	$(CURDIR)/../scripts/port-forward.sh --config=$(HELM_FILE)

.PHONY: stop-forward
stop-forward:
	-pkill -f "port-forward"
