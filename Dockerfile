ARG BASE_IMAGE_BUILDER=golang
ARG ALPINE_VERSION=3.23
ARG GO_VERSION=1.25

FROM ${BASE_IMAGE_BUILDER}:${GO_VERSION}-alpine${ALPINE_VERSION} AS builder
WORKDIR /tmp/gobuild
COPY . .
ENV CGO_ENABLED=0 \
    DISABLE_WARN_OUTSIDE_CONTAINER=1
RUN cd ./cmd/restAPI && go build .

FROM scratch
COPY --from=builder /tmp/gobuild/cmd/restAPI/restAPI /bin/restAPI
ENTRYPOINT [ "/bin/restAPI" ]
