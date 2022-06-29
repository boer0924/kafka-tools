FROM golang:1.17-alpine AS builder1
WORKDIR /app
COPY producer-random/go.mod .
COPY producer-random/go.sum .
COPY producer-random/vendor ./vendor
COPY producer-random/main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -mod vendor -o kpro main.go

FROM golang:1.17-alpine AS builder2
WORKDIR /app
COPY consumer-logger/go.mod .
COPY consumer-logger/go.sum .
COPY consumer-logger/vendor ./vendor
COPY consumer-logger/main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -mod vendor -o kcom main.go

FROM edenhill/kcat:1.7.1 as builder3

FROM golang:1.17-alpine AS builder4
WORKDIR /app
COPY topic/go.mod .
COPY topic/go.sum .
COPY topic/vendor ./vendor
COPY topic/main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -mod vendor -o ktop main.go

FROM alpine:3.10
ENV RUN_DEPS libcurl lz4-libs zstd-libs ca-certificates bash
RUN echo Installing ; \
#   apk add --no-cache --virtual .dev_pkgs $BUILD_DEPS $BUILD_DEPS_EXTRA && \
  apk add --no-cache $RUN_DEPS $RUN_DEPS_EXTRA && \
  rm -rf /usr/src/kcat && \
#   apk del .dev_pkgs && \
  rm -rf /var/cache/apk/*
COPY --from=builder1 /app/kpro /usr/bin/kpro
COPY --from=builder2 /app/kcom /usr/bin/kcom
COPY --from=builder3 /usr/bin/kcat /usr/bin/kcat
COPY --from=builder4 /app/ktop /usr/bin/ktop