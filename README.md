# HetznerCloud Exporter

[![Current Tag](https://img.shields.io/github/v/tag/promhippie/hcloud_exporter?sort=semver)](https://github.com/promhippie/hcloud_exporter) [![Build Status](https://github.com/promhippie/hcloud_exporter/actions/workflows/general.yml/badge.svg)](https://github.com/promhippie/hcloud_exporter/actions) [![Join the Matrix chat at https://matrix.to/#/#webhippie:matrix.org](https://img.shields.io/badge/matrix-%23webhippie-7bc9a4.svg)](https://matrix.to/#/#webhippie:matrix.org) [![Docker Size](https://img.shields.io/docker/image-size/promhippie/hcloud-exporter/latest)](https://hub.docker.com/r/promhippie/hcloud-exporter) [![Docker Pulls](https://img.shields.io/docker/pulls/promhippie/hcloud-exporter)](https://hub.docker.com/r/promhippie/hcloud-exporter) [![Go Reference](https://pkg.go.dev/badge/github.com/promhippie/hcloud_exporter.svg)](https://pkg.go.dev/github.com/promhippie/hcloud_exporter) [![Go Report Card](https://goreportcard.com/badge/github.com/promhippie/hcloud_exporter)](https://goreportcard.com/report/github.com/promhippie/hcloud_exporter) [![Codacy Badge](https://app.codacy.com/project/badge/Grade/0621f7fa70104074ad05ab9ac304d185)](https://www.codacy.com/gh/promhippie/hcloud_exporter/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=promhippie/hcloud_exporter&amp;utm_campaign=Badge_Grade)

An exporter for [Prometheus](https://prometheus.io/) that collects metrics from [Hetzner Cloud](https://console.hetzner.cloud).

## Install

You can download prebuilt binaries from our [GitHub releases](https://github.com/promhippie/hcloud_exporter/releases), or you can use our Docker images published on [Docker Hub](https://hub.docker.com/r/promhippie/hcloud-exporter/tags/) or [Quay](https://quay.io/repository/promhippie/hcloud-exporter?tab=tags). If you need further guidance how to install this take a look at our [documentation](https://promhippie.github.io/hcloud_exporter/#getting-started).

## Development

Make sure you have a working Go environment, for further reference or a guide take a look at the [install instructions](http://golang.org/doc/install.html). This project requires Go >= v1.11.

```bash
git clone https://github.com/promhippie/hcloud_exporter.git
cd hcloud_exporter

make generate build

./bin/hcloud_exporter -h
```

## Security

If you find a security issue please contact [thomas@webhippie.de](mailto:thomas@webhippie.de) first.

## Contributing

Fork -> Patch -> Push -> Pull Request

## Authors

-   [Thomas Boerger](https://github.com/tboerger)

## License

Apache-2.0

## Copyright

```console
Copyright (c) 2018 Thomas Boerger <thomas@webhippie.de>
```
