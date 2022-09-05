PROCESS = persistencecore
IMAGE_NAME = persistenceone/$(PROCESS)
TAG_NAME = latest
CONTAINER_NAME = container-$(PROCESS)
FILE = Dockerfile

DOCKER := $(shell which docker)

_DOCKER_ENV=\
	-e CHAIN_ID=pstake \
	-e PSTAKE_DIR=/pstaked \
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
	docker buildx build . --platform=linux/amd64 -f docker/$(PROCESS)/$(FILE) -t $(IMAGE_NAME):$(TAG_NAME)

docker-run:
	docker run --rm -it --name=$(CONTAINER_NAME) \
		$(_DOCKER_VOLUME) $(_DOCKER_ARGS) $(_DOCKER_ENV) \
		$(IMAGE_NAME):$(TAG_NAME) /bin/bash

docker-push:
	docker push $(IMAGE_NAME):$(TAG_NAME)

docker-build-push: docker-build docker-push

docker-clean-container:
	-docker stop $(CONTAINER_NAME)
	-docker rm $(CONTAINER_NAME)

docker-clean: docker-clean-container
	-docker rmi $(IMAGE_NAME):$(TAG_NAME)