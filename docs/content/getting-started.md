---
title: "Getting Started"
date: 2018-05-02T00:00:00+00:00
anchor: "getting-started"
weight: 10
---

## Installation

We won't cover further details how to properly setup [Prometheus](https://prometheus.io) itself, we will only cover some basic setup based on [docker-compose](https://docs.docker.com/compose/). But if you want to run this exporter without [docker-compose](https://docs.docker.com/compose/) you should be able to adopt that to your needs.

First of all we need to prepare a configuration for [Prometheus](https://prometheus.io) that includes the exporter as a target based on a static host mapping which is just the [docker-compose](https://docs.docker.com/compose/) container name, e.g. `hcloud-exporter`.

{{< highlight txt >}}
global:
  scrape_interval: 1m
  scrape_timeout: 10s
  evaluation_interval: 1m

scrape_configs:
- job_name: hcloud
  static_configs:
  - targets:
    - hcloud-exporter:9501
{{< / highlight >}}

After preparing the configuration we need to create the `docker-compose.yml` within the same folder, this `docker-compose.yml` starts a simple [Prometheus](https://prometheus.io) instance together with the exporter. Don't forget to update the exporter envrionment variables with the required credentials.

{{< highlight txt >}}
version: '2'

volumes:
  prometheus:

services:
  prometheus:
    image: prom/prometheus:v2.4.3
    restart: always
    ports:
      - 9090:9090
    volumes:
      - prometheus:/prometheus
      - ./prometheus.yml:/etc/prometheus/prometheus.yml

  hcloud-exporter:
    image: promhippie/hcloud-exporter:latest
    restart: always
    environment:
      - HCLOUD_EXPORTER_LOG_PRETTY=true
      - HCLOUD_EXPORTER_TOKEN=uDgX6TkZVGx7c94jPAff5cfJdym9MLekiveDgN7Oq5dyOXxl4Uu9qkpcC1muILGW
{{< / highlight >}}

Since our `latest` Docker tag always refers to the `master` branch of the Git repository you should always use some fixed version. You can see all available tags at our [DockerHub repository](https://hub.docker.com/r/promhippie/hcloud-exporter/tags/), there you will see that we also provide a manifest, you can easily start the exporter on various architectures without any change to the image name. You should apply a change like this to the `docker-compose.yml`:

{{< highlight diff >}}
  hcloud-exporter:
-   image: promhippie/hcloud-exporter:latest
+   image: promhippie/hcloud-exporter:0.1.0
    restart: always
    environment:
      - HCLOUD_EXPORTER_LOG_PRETTY=true
      - HCLOUD_EXPORTER_TOKEN=uDgX6TkZVGx7c94jPAff5cfJdym9MLekiveDgN7Oq5dyOXxl4Uu9qkpcC1muILGW
{{< / highlight >}}

If you want to access the exporter directly you should bind it to a local port, otherwise only [Prometheus](https://prometheus.io) will have access to the exporter. For debugging purpose or just to discover all available metrics directly you can apply this change to your `docker-compose.yml`, after that you can access it directly at [http://localhost:9501/metrics](http://localhost:9501/metrics):

{{< highlight diff >}}
  hcloud-exporter:
    image: promhippie/hcloud-exporter:latest
    restart: always
+   ports:
+     - 127.0.0.1:9501:9501
    environment:
      - HCLOUD_EXPORTER_LOG_PRETTY=true
      - HCLOUD_EXPORTER_TOKEN=uDgX6TkZVGx7c94jPAff5cfJdym9MLekiveDgN7Oq5dyOXxl4Uu9qkpcC1muILGW
{{< / highlight >}}

That's all, the exporter should be up and running. Have fun with it and hopefully you will gather interesting metrics and never run into issues.

## Kubernetes

Currently we have not prepared a deployment for Kubernetes, but this is something we will provide for sure. Most interesting will be the integration into the [Prometheus Operator](https://coreos.com/operators/prometheus/docs/latest/), so stay tuned.

## Configuration

HCLOUD_EXPORTER_TOKEN
: Access token for the HetznerCloud API, required for authentication

HCLOUD_EXPORTER_LOG_LEVEL
: Only log messages with given severity, defaults to `info`

HCLOUD_EXPORTER_LOG_PRETTY
: Enable pretty messages for logging, defaults to `false`

HCLOUD_EXPORTER_WEB_ADDRESS
: Address to bind the metrics server, defaults to `0.0.0.0:9501`

HCLOUD_EXPORTER_WEB_PATH
: Path to bind the metrics server, defaults to `/metrics`

HCLOUD_EXPORTER_REQUEST_TIMEOUT
: Request timeout as duration, defaults to `5s`

HCLOUD_EXPORTER_COLLECTOR_FLOATING_IPS
: Enable collector for floating IPs, defaults to  `true`

HCLOUD_EXPORTER_COLLECTOR_IMAGES
: Enable collector for images, defaults to `true`

HCLOUD_EXPORTER_COLLECTOR_SERVERS
: Enable collector for servers, defaults to `true`

HCLOUD_EXPORTER_COLLECTOR_SSH_KEYS
: Enable collector for SSH keys, defaults to `true`

## Metrics

hcloud_request_duration_seconds
: Histogram of latencies for requests to the HetznerCloud API per collector

hcloud_request_failures_total
: Total number of failed requests to the HetznerCloud API per collector

hcloud_floating_ip_active
: If 1 the floating IP is used by a server, 0 otherwise

hcloud_image_active
: If 1 the image is used by a server, 0 otherwise

hcloud_image_size_bytes
: Size of the image in bytes

hcloud_image_disk_bytes
: Size if the disk for the image in bytes

hcloud_image_created_timestamp
: Timestamp when the image have been created

hcloud_image_deprecated_timestamp
: Timestamp when the image will be deprecated, 0 if not deprecated

hcloud_server_running
: If 1 the server is running, 0 otherwise

hcloud_server_created_timestamp
: Timestamp when the server have been created

hcloud_server_included_traffic_bytes
: Included traffic for the server in bytes

hcloud_server_outgoing_traffic_bytes
: Outgoing traffic from the server in bytes

hcloud_server_incming_traffic_bytes
: Ingoing traffic to the server in bytes

hcloud_server_cores
: Server number of cores

hcloud_server_memory_bytes
: Server memory in bytes

hcloud_server_disk_bytes
: Server disk in bytes

hcloud_server_price_hourly
: Price of the server billed hourly in €

hcloud_server_price_monthly
: Price of the server billed monthly in €

hcloud_ssh_key
: Information about SSH keys in your HetznerCloud project
