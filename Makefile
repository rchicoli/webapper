NAME = rchicoli/app
VERSION = 0.0.1-dev
WORKDIR = /go/src/app
BINARY = app

.PHONY: all build tag release

all: binary build

binary:
	docker run --rm -ti -v $(PWD):$(WORKDIR) -w $(WORKDIR) golang:1.7.1-alpine go build -v

build:
	docker build --rm -t $(NAME):$(VERSION) .

tag:
	docker tag $(NAME):$(VERSION) $(NAME):latest

release: tag
        docker push $(NAME):$(VERSION)
        docker push $(NAME):latest
