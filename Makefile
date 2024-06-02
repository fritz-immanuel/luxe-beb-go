# Copyright 2018 Shopfazz Authors.

CI_CONTRACT := docker-1

# Which architecture to build
ARCH ?= amd64

# The docker image name
IMAGE := luxe-beb-go

# The docker registry ip
REGISTRY := atomicindonesia
#REGISTRY := rickyhandoko

# This version-strategy uses git tags to set the version string
VERSION := $(shell git describe --tags --always --dirty)

.PHONY: migrate
migrate:
	go run -v ./main.go

.PHONY: test
test:
	go test -v -cover -p 1 ./... -tags test 

.PHONY: build
build:
	CGO_ENABLED=0 GOARCH=${ARCH} go install ./cmd/...

.PHONY: build-linux
build-linux:
	GOOS=linux CGO_ENABLED=0 GOARCH=${ARCH} go install ./cmd/...

.PHONY: ci-docker-image
ci-docker-image:
	echo "building the $(IMAGE) container..."
	docker build --build-arg TOKEN="$(GITHUB_TOKEN)" --label "version=$(VERSION)" -t $(CI_DOCKER_IMAGE_OUTPUT) .


docker: Dockerfile
	echo "building the $(IMAGE) container..."
	docker build --label "version=$(VERSION)" -t $(IMAGE):$(VERSION) .

push-docker: .push-docker
.push-docker:
	docker tag $(IMAGE):$(VERSION) $(REGISTRY)/$(IMAGE):$(VERSION)
	docker push $(REGISTRY)/$(IMAGE):$(VERSION)
	echo "pushed: $(REGISTRY)/$(IMAGE):$(VERSION)"