
###############################################################################
###                              Docker commands                            ###
###############################################################################

# Docker variables
DOCKER := $(shell which docker)
DOCKER_COMPOSE := $(shell which docker-compose)

all: docker-compose

docker-compose: docker-compose-clean
	$(DOCKER_COMPOSE) up

docker-compose-clean:
	$(DOCKER_COMPOSE) rm --force --stop

docker-exec-relayer:
	$(DOCKER_COMPOSE) exec relayer /bin/bash

###############################################################################
###                      Docker support build commands                      ###
###############################################################################

GAIA_TAG_NAME := v7.0.2
HERMES_TAG_NAME := 1.0.0-rc.1

docker-build-gaiad:
	$(DOCKER) buildx build ${DOCKER_ARGS} \
		--build-arg GAIA_TAG_NAME=$(GAIA_TAG_NAME) \
		-t persistenceone/gaiad:$(GAIA_TAG_NAME) \
		. -f gaiad.Dockerfile

docker-build-push-gaiad: docker-build-gaiad
	$(DOCKER) push persistenceone/gaiad:$(GAIA_TAG_NAME)

docker-build-hermes:
	$(DOCKER) buildx build ${DOCKER_ARGS} \
		-t persistenceone/hermes:$(HERMES_TAG_NAME) \
		. -f hermes.Dockerfile

docker-build-push-hermes: docker-build-hermes
	$(DOCKER) push persistenceone/hermes:$(HERMES_TAG_NAME)