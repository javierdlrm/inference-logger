# Image URL to use all building/pushing image targets
IMG ?= javierdlrm/inferencelogger:latest

docker: docker-build docker-push

# Build the docker image
docker-build:
	docker build . -t ${IMG}

# Push the docker image
docker-push:
	docker push ${IMG}