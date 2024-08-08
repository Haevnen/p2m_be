FROM golang:1.22.5-alpine3.19
RUN apk --no-cache add git g++ openssh-client

ENV OAPI_VERSION=v2.3.0
ENV MOCKGEN_VERSION=v0.4.0
ENV STRINGER_VERSION=v0.22.0
ENV GO_IMPORTS_VERSION=v0.22.0

RUN go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@${OAPI_VERSION}
RUN go install golang.org/x/tools/cmd/stringer@${STRINGER_VERSION}
RUN go install go.uber.org/mock/mockgen@${MOCKGEN_VERSION}
RUN go install golang.org/x/tools/cmd/goimports@${GO_IMPORTS_VERSION}
