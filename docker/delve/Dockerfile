FROM golang:1.22.5-alpine3.19

ENV AIR_VERSION=v1.52.2
ENV DELVE_VERSION=v1.22.1

RUN go install github.com/air-verse/air@"${AIR_VERSION}"
RUN go install github.com/go-delve/delve/cmd/dlv@"${DELVE_VERSION}"
