# STEP 1 - build
FROM golang:1.12.7-alpine3.10 as build

ARG GOOS=linux
ARG GOARCH=amd64
ARG GOARM=

ARG APP_NAME=webapper
ARG APP_USER=app
ARG APP_GROUP=app
ARG APP_HOME=/home/app

RUN apk add --update curl git

WORKDIR $GOPATH/src/github.com/rchicoli/webapper
COPY . .

RUN go get -d -v ./...

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -a -installsuffix cgo -o ${APP_NAME} main.go

# STEP 2 - release
FROM alpine:3.10

ARG GOOS=linux
ARG GOARCH=amd64
ARG GOARM=

ARG APP_NAME=webapper
ARG APP_USER=app
ARG APP_GROUP=app
ARG APP_HOME=/home/app

WORKDIR ${APP_HOME}

RUN apk --no-cache add ca-certificates shadow

RUN groupadd -r ${APP_GROUP} && useradd --no-log-init -r -g ${APP_GROUP} -d ${APP_HOME} -m -s /sbin/nologin ${APP_USER}

COPY --from=build  /go/src/github.com/rchicoli/webapper/${APP_NAME} ${APP_HOME}/

USER ${APP_USER}

EXPOSE 8080

CMD ["/home/app/webapper"]
