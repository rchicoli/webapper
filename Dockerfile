FROM alpine:3.4

COPY app /app

RUN apk add --update curl

HEALTHCHECK --interval=10s --timeout=3s CMD curl --fail http://localhost:8080/health || exit 1

CMD ["/app"]
