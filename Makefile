# Copyright 2020 Keyporttech Inc.
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#     http://www.apache.org/licenses/LICENSE-2.0
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

DOCKER_REGISTRY="registry.keyporttech.com"
DOCKERHUB_REGISTRY="keyporttech"
IMAGE=gitea-webhook-interceptor
VERSION = $(shell cat main.go | grep "Version = " | sed 's/var Version = //g' | tr -d '"')


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
