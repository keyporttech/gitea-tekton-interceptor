
DOCKER_REGISTRY="registry.keyporttech.com:30243"
DOCKERHUB_REGISTRY="keyporttech"
IMAGE=gitea-webbhook-interceptor
VERSION = $(shell git describe --abbrev=0)


compile:
	@echo "making..."

	CGO_ENABLED=0 go build -a \
		-o ./ \
		--installsuffix cgo \
		--ldflags="--s" \
		.
	@echo "OK"

.PHONY: build

test:
	@echo "testing..."
	go test
	@echo "OK"
.PHONY: test

build: compile test
	
.PHONY: build

docker:
	docker build ./ -t ${DOCKER_REGISTRY}/${IMAGE}:${VERSION}
	docker build ./ -t ${DOCKERHUB_REGISTRY}/${IMAGE}:${VERSION}
.PHONY: docker

docker-publish: docker
	docker push ${DOCKER_REGISTRY}/${IMAGE}:${VERSION}
	docker push ${DOCKERHUB_REGISTRY}/${IMAGE}:${VERSION}
.PHONY: docker-publish
