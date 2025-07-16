FROM --platform=$BUILDPLATFORM golang:1.24.5-alpine3.21@sha256:42b91b1364b79204d9b421e03487efbe73d267211d8e24d8755234ba578451e4 AS builder

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

FROM alpine:3.22@sha256:4bcff63911fcb4448bd4fdacec207030997caf25e9bea4045fa6c8c44de311d1

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
