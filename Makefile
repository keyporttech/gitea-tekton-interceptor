
DOCKER_REGISTRY="registry.keyporttech.com:30243"
DOCKERHUB_REGISTRY="keyporttech"
IMAGE=gitea-webbhook-interceptor
VERSION = $(shell git describe --exact-match --abbrev=0)


build:
	@echo "making..."

	CGO_ENABLED=0 go build -a \
		-o ./ \
		--installsuffix cgo \
		--ldflags="--s" \
		.
	@echo "OK"
.PHONY: build

docker:
	docker build ./ -t ${DOCKER_REGISTRY}/${IMAGE}:${VERSION}
	docker build ./ -t ${DOCKERHUB_REGISTRY}/${IMAGE}:${VERSION}
.PHONY: docker

docker-publish: docker
	docker push ${DOCKER_REGISTRY}/${IMAGE}:${VERSION}
	docker push ${DOCKERHUB_REGISTRY}/${IMAGE}:${VERSION}
.PHONY: docker-publish
