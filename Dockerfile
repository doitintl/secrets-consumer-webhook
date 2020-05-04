ARG GO_VERSION=1.14

FROM golang:${GO_VERSION}-alpine AS builder
ARG VERSION
ARG COMMIT

RUN apk add --update --no-cache ca-certificates make git curl

RUN mkdir -p /build
WORKDIR /build

ARG GOPROXY

RUN mkdir -p /build
WORKDIR /build

COPY go.* /build/
RUN go mod download

COPY . /build
RUN go build -ldflags="-X github.com/doitintl/secrets-consumer-webhook/version.version=${VERSION} -X github.com/doitintl/secrets-consumer-webhook/version.gitCommitID=${COMMIT}"
RUN cp secrets-consumer-webhook /usr/local/bin/
RUN chmod a+x /usr/local/bin/secrets-consumer-webhook

FROM alpine

RUN apk add --update libcap && rm -rf /var/cache/apk/*

COPY --from=builder /usr/local/bin/secrets-consumer-webhook /usr/local/bin/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENV DEBUG false
USER 65534

ENTRYPOINT ["/usr/local/bin/secrets-consumer-webhook"]
