FROM alpine:3.4

# ARG GOOS=linux
# ARG GOARCH=amd64
# ARG GOARM=

ARG APP_NAME
ENV APP_NAME ${APP_NAME:-webapper}

COPY $APP_NAME /usr/bin/$APP_NAME

RUN apk add --update curl

HEALTHCHECK --interval=10s --timeout=3s CMD curl --silent --fail http://localhost:8080/health || exit 1

CMD ["webapper"]
