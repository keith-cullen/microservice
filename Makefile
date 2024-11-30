REG = localhost:5000
IMG = microservice
TAG = latest

#.PHONY: docker
#docker:	docker-build docker-push

.PHONY: all
all:
	go build

.PHONY: docker-build
docker-build:
	docker build --tag $(REG)/$(IMG):$(TAG) .

.PHONY: docker-push
docker-push:
	docker push $(REG)/$(IMG):$(TAG)
	curl http://$(REG)/v2/$(IMG)/tags/list
