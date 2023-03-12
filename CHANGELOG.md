# Changelog for 1.2.2

The following sections list the changes for 1.2.2.

## Summary

 * Fix #72: Fix index out of range issue within server metrics
 * Fix #74: Another fix for go routines within server metrics

## Details

 * Bugfix #72: Fix index out of range issue within server metrics

   The code has not checked if an index have been really available within the server metrics API
   response. With this fix it gets properly handled.

   https://github.com/promhippie/hcloud_exporter/issues/72

 * Bugfix #74: Another fix for go routines within server metrics

   We disabled the server metrics by default for now until the implementation is really stable to
   avoid any side effects. I have reintroduced routines, otherwise the scrapetime will be far too
   high. This time I used wait groups to get everything handled properly.

   https://github.com/promhippie/hcloud_exporter/issues/74


# Changelog for 1.2.1

The following sections list the changes for 1.2.1.

## Summary

 * Fix #70: Fix go routine errors within server metrics

## Details

 * Bugfix #70: Fix go routine errors within server metrics

   We fixed a go routines issue within the new server metrics. We just got rid of the routines to
   avoid any errors related to sending to closed channels.

   https://github.com/promhippie/hcloud_exporter/issues/70


# Changelog for 1.2.0

The following sections list the changes for 1.2.0.

## Summary

 * Chg #53: Integrate standard web config
 * Chg #67: Add collector for server metrics

## Details

 * Change #53: Integrate standard web config

   We integrated the new web config from the Prometheus toolkit which provides a configuration
   for TLS support and also some basic builtin authentication. For the detailed configuration
   you check out the documentation.

   https://github.com/promhippie/hcloud_exporter/issues/53

 * Change #67: Add collector for server metrics

   Hetzner Cloud collects basic metrics on the hypervisor-level for each server. We have added a
   new collector which scrapes the latest available metric point for each running server. It is
   enabled by default.

   https://github.com/promhippie/hcloud_exporter/pull/67


# Changelog for 1.1.0

The following sections list the changes for 1.1.0.

## Summary

 * Chg #21: Add collector for volumes
 * Chg #24: Refactor build tools and project structure
 * Chg #25: Drop darwin/386 release builds
 * Chg #39: Add collector for load balancers

## Details

 * Change #21: Add collector for volumes

   We have added a new optional collector, which is disabled by default, to gather metrics about
   the volumes part of the configured Hetzner Cloud project.

   https://github.com/promhippie/hcloud_exporter/issues/21

 * Change #24: Refactor build tools and project structure

   To have a unified project structure and build tooling we have integrated the same structure we
   already got within our GitHub exporter.

   https://github.com/promhippie/hcloud_exporter/issues/24

 * Change #25: Drop darwin/386 release builds

   We dropped the build of 386 builds on Darwin as this architecture is not supported by current Go
   versions anymore.

   https://github.com/promhippie/hcloud_exporter/issues/25

 * Change #39: Add collector for load balancers

   We have added a new optional collector, which is enabled by default, to gather metrics about all
   loadbalancers part of the configured Hetzner Cloud project.

   https://github.com/promhippie/hcloud_exporter/issues/39


# Changelog for 1.0.0

The following sections list the changes for 1.0.0.

## Summary

 * Chg #19: Add `vat` labels for net and gross values

## Details

 * Change #19: Add `vat` labels for net and gross values

   Added a new `vat` label for `gross` or `net` values to the `hcloud_server_price` metric.
   Depending on the setup this can be a breaking change and it may be necessary to adjust some
   dashboards and alerting rules.

   https://github.com/promhippie/hcloud_exporter/pull/19


# Changelog for 0.2.0

The following sections list the changes for 0.2.0.

## Summary

 * Chg #17: Add pricing collector
 * Chg #18: Add new metric to see if backups enabled

## Details

 * Change #17: Add pricing collector

   We added a new collector to gather information about the pricings, that way somebody could do
   calculations how much the costs are increasing or decreasing by sclae up or sclae down. The new
   collector includes new metrics named `hcloud_pricing_floating_ip`,
   `hcloud_pricing_image`, `hcloud_pricing_server_backup` and
   `hcloud_pricing_traffic`.

   https://github.com/promhippie/hcloud_exporter/pull/17

 * Change #18: Add new metric to see if backups enabled

   We added a new metric named `hcloud_server_backup` which indicates if a server got backups
   enabled or not, that way somebody could add some alerting if a server is missing a backup.

   https://github.com/promhippie/hcloud_exporter/pull/18


# Changelog for 0.1.1

The following sections list the changes for 0.1.1.

## Summary

 * Fix #11: Fix typo within `hcloud_server_incoming_traffic_bytes`
 * Chg #13: Pin go version to 1.10

## Details

 * Bugfix #11: Fix typo within `hcloud_server_incoming_traffic_bytes`

   We fixed a typo within the `hcloud_server_incoming_traffic_bytes` metric where we were just
   missing a tiny single letter.

   https://github.com/promhippie/hcloud_exporter/pull/11

 * Change #13: Pin go version to 1.10

   To make sure we got something nearly like reproducible builds and to fix the builds we should pin
   the build dependencies like the Go version to make sure it is always buildable.

   https://github.com/promhippie/hcloud_exporter/pull/13


# Changelog for 0.1.0

The following sections list the changes for 0.1.0.

## Summary

 * Chg #23: Initial release of basic version

## Details

 * Change #23: Initial release of basic version

   Just prepared an initial basic version which could be released to the public.

   https://github.com/promhippie/hcloud_exporter/issues/23


