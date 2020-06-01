# Version of the image
VERSION ?= v1beta1
# Image URL to use all building/pushing image targets
IMG ?= javierdlrm/inference-logger:${VERSION}

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# Build and publish as docker image
docker: fmt vet docker-build docker-push

# Build the docker image
docker-build:
	docker build . -t ${IMG}

# Push the docker image
docker-push:
	docker push ${IMG}