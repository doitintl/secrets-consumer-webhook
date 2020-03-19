ARG GO_VERSION=1.14

FROM golang:${GO_VERSION}-alpine AS builder

RUN apk add --update --no-cache ca-certificates make git curl

RUN mkdir -p /build
WORKDIR /build

ARG GOPROXY

RUN mkdir -p /build
WORKDIR /build

COPY go.* /build/
RUN go mod download

COPY . /build
RUN go install .

FROM alpine:3.10

RUN apk add --update libcap && rm -rf /var/cache/apk/*

COPY --from=builder /go/bin/secrets-consumer-webhook /usr/local/bin/secrets-consumer-webhook
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENV DEBUG false
USER 65534

ENTRYPOINT ["/usr/local/bin/secrets-consumer-webhook"]
