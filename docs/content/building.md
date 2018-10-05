---
title: "Building"
date: 2018-05-02T00:00:00+00:00
anchor: "building"
weight: 20
---

As this project is built with Go you need to install Go first. The installation of Go is out of the scope of this document, please follow the [official documentation](https://golang.org/doc/install). After the installation of Go you need to get the sources:

{{< highlight txt >}}
go get -d github.com/promhippie/hcloud_exporter
cd $GOPATH/src/github.com/promhippie/hcloud_exporter
{{< / highlight >}}

All required tool besides Go itself are bundled or getting automatically installed within the `GOPATH`. We are using [retool](https://github.com/twitchtv/retool) to keep the used tools consistent and [dep](https://github.com/golang/dep) to manage the dependencies. All commands to build this project are part of our `Makefile`.

{{< highlight txt >}}
# install retool
make retool

# sync dependencies
make sync

# generate code
make generate

# build binary
make build
{{< / highlight >}}

Finally you should have the binary within the `bin/` folder now, give it a try with `./bin/hcloud_exporter -h` to see all available options.
