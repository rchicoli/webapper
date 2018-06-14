
APP_OWNER     = rchicoli
APP_NAME      = webapper
APP_VERSION   ?= latest

DOCKER_IMAGE  ?= $(APP_OWNER)/$(APP_NAME)
WORKDIR		  ?= /go/src/app

.PHONY: all binary build clean release run tag

all: test binary build

binary:
	docker run --rm -v $(PWD):$(WORKDIR) -w $(WORKDIR) golang:1.7.1-alpine go build -ldflags '-extldflags "-static"' -o $(APP_NAME) main.go

build:
	docker build --rm --build-arg APP_NAME=$(APP_NAME) -t $(DOCKER_IMAGE):$(APP_VERSION) .

run:
	docker run --rm $(DOCKER_IMAGE):$(APP_VERSION)

attach:
	docker run --rm -ti $(DOCKER_IMAGE):$(APP_VERSION) /bin/sh

push:
	docker push $(DOCKER_IMAGE):$(APP_VERSION)

tag:
	docker tag $(DOCKER_IMAGE):$(APP_VERSION) $(DOCKER_IMAGE):latest

test:
	go test -v

release: tag
	docker push $(DOCKER_IMAGE):$(APP_VERSION)
	docker push $(DOCKER_IMAGE):latest

clean:
	rm -f $(APP_NAME)
	docker rmi $(DOCKER_IMAGE):$(APP_VERSION)
