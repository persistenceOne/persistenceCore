PROCESS = persistencecore
DOCKER_IMAGE_NAME = persistenceone/$(PROCESS)
DOCKER_TAG_NAME = v6.0.0
CONTAINER_NAME = container-$(PROCESS)
DOCKER_FILE = Dockerfile

DOCKER := $(shell which docker)

_DOCKER_ENV=\
	-e CHAIN_ID=persistencecore \
	-e COINS="100000000000000000stake" \
	$(DOCKER_ENV)
_DOCKER_VOLUME=\
	-v /tmp:/tmp \
	$(DOCKER_VOLUME)
_DOCKER_ARGS=\
	-p 26656:26656 \
	-p 26657:26657 \
	-p 8080:8080 \
	--network=host \
	$(DOCKER_ARGS)

# Command for dockerization
docker-build:
	$(DOCKER) buildx build . \
		$(DOCKER_BUILD_ARGS) \
		-f docker/$(PROCESS)/$(DOCKER_FILE) \
		-t $(DOCKER_IMAGE_NAME):$(DOCKER_TAG_NAME)

docker-run:
	$(DOCKER) run --rm -it --name=$(CONTAINER_NAME) \
		$(_DOCKER_VOLUME) $(_DOCKER_ARGS) $(_DOCKER_ENV) \
		$(DOCKER_IMAGE_NAME):$(DOCKER_TAG_NAME) /bin/bash

docker-push:
	$(DOCKER) push $(DOCKER_IMAGE_NAME):$(DOCKER_TAG_NAME)

docker-build-push: docker-build docker-push

docker-clean-container:
	-$(DOCKER) stop $(CONTAINER_NAME)
	-$(DOCKER) rm $(CONTAINER_NAME)

docker-clean: docker-clean-container
	-$(DOCKER) rmi $(DOCKER_IMAGE_NAME):$(DOCKER_TAG_NAME)
