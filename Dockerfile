FROM alpine:3.4

ARG APP_NAME
ENV APP_NAME ${APP_NAME:-app}

COPY $APP_NAME /$APP_NAME

RUN apk add --update curl

HEALTHCHECK --interval=10s --timeout=3s CMD curl --fail http://localhost:8080/health || exit 1

CMD ["$APP_NAME"]
