# HetznerCloud Exporter

[![Build Status](http://github.dronehippie.de/api/badges/promhippie/hcloud_exporter/status.svg)](http://github.dronehippie.de/promhippie/hcloud_exporter)
[![Stories in Ready](https://badge.waffle.io/promhippie/hcloud_exporter.svg?label=ready&title=Ready)](http://waffle.io/promhippie/hcloud_exporter)
[![Join the Matrix chat at https://matrix.to/#/#webhippie:matrix.org](https://img.shields.io/badge/matrix-%23webhippie-7bc9a4.svg)](https://matrix.to/#/#webhippie:matrix.org)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/92400592c23c4ea4bd3cf2448eb110ef)](https://www.codacy.com/app/promhippie/hcloud_exporter?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=promhippie/hcloud_exporter&amp;utm_campaign=Badge_Grade)
[![Go Doc](https://godoc.org/github.com/promhippie/hcloud_exporter?status.svg)](http://godoc.org/github.com/promhippie/hcloud_exporter)
[![Go Report](http://goreportcard.com/badge/github.com/promhippie/hcloud_exporter)](http://goreportcard.com/report/github.com/promhippie/hcloud_exporter)
[![](https://images.microbadger.com/badges/image/promhippie/hcloud-exporter.svg)](http://microbadger.com/images/promhippie/hcloud-exporter "Get your own image badge on microbadger.com")

An exporter for [Prometheus](https://prometheus.io/) that collects metrics from [Hetzner Cloud](https://console.hetzner.cloud).

## Install

You can download prebuilt binaries from our [GitHub releases](https://github.com/promhippie/hcloud_exporter/releases), or you can use our Docker images published on [Docker Hub](https://hub.docker.com/r/promhippie/hcloud_exporter/tags/). If you need further guidance how to install this take a look at our [documentation](https://promhippie.github.io/hcloud_exporter/#getting-started).

## Development

Make sure you have a working Go environment, for further reference or a guide take a look at the [install instructions](http://golang.org/doc/install.html). This project requires Go >= v1.8.

```bash
go get -d github.com/promhippie/hcloud_exporter
cd $GOPATH/src/github.com/promhippie/hcloud_exporter

# install retool
make retool

# sync dependencies
make sync

# generate code
make generate

# build binary
make build

./bin/hcloud_exporter -h
```

## Security

If you find a security issue please contact thomas@webhippie.de first.

## Contributing

Fork -> Patch -> Push -> Pull Request

## Authors

* [Thomas Boerger](https://github.com/tboerger)

## License

Apache-2.0

## Copyright

```console
Copyright (c) 2018 Thomas Boerger <thomas@webhippie.de>
```
