FROM --platform=$BUILDPLATFORM golang:1.27.3-alpine@sha256:dbaa054703263d06c8b49576667418473b6231f8c0fa69dee103207fdb8e3654 AS builder

RUN apk add --no-cache -U git curl
RUN sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b /usr/local/bin

WORKDIR /go/src/exporter
COPY . /go/src/exporter/

RUN --mount=type=cache,target=/go/pkg \
    go mod download -x

ARG TARGETOS
ARG TARGETARCH

RUN --mount=type=cache,target=/go/pkg \
    --mount=type=cache,target=/root/.cache/go-build \
    task generate build GOOS=${TARGETOS} GOARCH=${TARGETARCH}

FROM alpine:3.24@sha256:ce9e1e6d23f04a0b7593445a062a9fc73ed52f8260b9396f77a958295ad18ee7

RUN apk add --no-cache ca-certificates mailcap && \
    addgroup -g 1337 exporter && \
    adduser -D -u 1337 -h /var/lib/exporter -G exporter exporter

EXPOSE 9501
VOLUME ["/var/lib/exporter"]
ENTRYPOINT ["/usr/bin/hcloud_exporter"]
HEALTHCHECK CMD ["/usr/bin/hcloud_exporter", "health"]

COPY --from=builder /go/src/exporter/bin/hcloud_exporter /usr/bin/hcloud_exporter
WORKDIR /var/lib/exporter
USER exporter
