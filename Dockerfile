FROM --platform=$BUILDPLATFORM golang:1.25.4-alpine3.21@sha256:1b84283ebeef726bc5fa9fec8deb36828aabfa12fe3a28b4fb0a4b2eafafe38c AS builder

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

FROM alpine:3.22@sha256:4b7ce07002c69e8f3d704a9c5d6fd3053be500b7f1c69fc0d80990c2ad8dd412

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
