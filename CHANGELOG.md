# Changelog for 2.1.0

The following sections list the changes for 2.1.0.

## Summary

 * Chg #258: Switch to official logging library
 * Chg #272: Add type to IP pricing and add metrics for primary IPs

## Details

 * Change #258: Switch to official logging library

   Since there have been a structured logger part of the Go standard library we
   thought it's time to replace the library with that. Be aware that log messages
   should change a little bit.

   https://github.com/promhippie/hcloud_exporter/issues/258

 * Change #272: Add type to IP pricing and add metrics for primary IPs

   Since the client SDK has deprecated the previous handling for the pricing of IP
   addresses we had to update the metrics to include the type and location of the
   IPs. Besides that we have also added metrics for the pricing of the primary IP
   addresses.

   https://github.com/promhippie/hcloud_exporter/pull/272


# Changelog for 2.0.0

The following sections list the changes for 2.0.0.

## Summary

 * Fix #246: Fetch metrics for all servers
 * Chg #240: Improve pricing error handling
 * Chg #240: New traffic pricing metrics because of deprecation

## Details

 * Bugfix #246: Fetch metrics for all servers

   For previous versions we have used the wrong client function to gather the list
   of servers for the server metrics, this hvae been fixed by using a function that
   automatically fetches all servers by iterating of the pagination.

   https://github.com/promhippie/hcloud_exporter/issues/246

 * Change #240: Improve pricing error handling

   So far we always existed the scraping if there have been any kind of error while
   parsing the metric values, from now on we are logging an error but continue to
   provide the remaining metrics to avoid loosing unrelated metrics.

   https://github.com/promhippie/hcloud_exporter/pull/240

 * Change #240: New traffic pricing metrics because of deprecation

   The previous traffic pricing metrics have been deprecated and got to be replaced
   by new metrics as the new metrics have been split between service type like load
   balancers and server types.

   https://github.com/promhippie/hcloud_exporter/issues/248
   https://github.com/promhippie/hcloud_exporter/pull/240


# Changelog for 1.3.0

The following sections list the changes for 1.3.0.

## Summary

 * Chg #193: Read secrets form files
 * Chg #193: Integrate standard web config
 * Enh #193: Integrate option pprof profiling

## Details

 * Change #193: Read secrets form files

   We have added proper support to load secrets like the password from files or
   from base64-encoded strings. Just provide the flags or environment variables for
   token or private key with a DSN formatted string like `file://path/to/file` or
   `base64://Zm9vYmFy`.

   https://github.com/promhippie/hcloud_exporter/pull/193

 * Change #193: Integrate standard web config

   We integrated the new web config from the Prometheus toolkit which provides a
   configuration for TLS support and also some basic builtin authentication. For
   the detailed configuration you can check out the documentation.

   https://github.com/promhippie/hcloud_exporter/pull/193

 * Enhancement #193: Integrate option pprof profiling

   We have added an option to enable a pprof endpoint for proper profiling support
   with the help of tools like Parca. The endpoint `/debug/pprof` can now
   optionally be enabled to get the profiling details for catching potential memory
   leaks.

   https://github.com/promhippie/hcloud_exporter/pull/193


# Changelog for 1.2.3

The following sections list the changes for 1.2.3.

## Summary

 * Fix #175: Correctly read loadbalancer traffic

## Details

 * Bugfix #175: Correctly read loadbalancer traffic

   We used a wrong attribute to read the loadbalancer traffic which resulted in
   missing metrics for the realtime traffic in and out for all loadbalancers. With
   this fix you should be able to use the metrics.

   https://github.com/promhippie/hcloud_exporter/issues/175


# Changelog for 1.2.2

The following sections list the changes for 1.2.2.

## Summary

 * Fix #72: Fix index out of range issue within server metrics
 * Fix #74: Another fix for go routines within server metrics

## Details

 * Bugfix #72: Fix index out of range issue within server metrics

   The code has not checked if an index have been really available within the
   server metrics API response. With this fix it gets properly handled.

   https://github.com/promhippie/hcloud_exporter/issues/72

 * Bugfix #74: Another fix for go routines within server metrics

   We disabled the server metrics by default for now until the implementation is
   really stable to avoid any side effects. I have reintroduced routines, otherwise
   the scrapetime will be far too high. This time I used wait groups to get
   everything handled properly.

   https://github.com/promhippie/hcloud_exporter/issues/74


# Changelog for 1.2.1

The following sections list the changes for 1.2.1.

## Summary

 * Fix #70: Fix go routine errors within server metrics

## Details

 * Bugfix #70: Fix go routine errors within server metrics

   We fixed a go routines issue within the new server metrics. We just got rid of
   the routines to avoid any errors related to sending to closed channels.

   https://github.com/promhippie/hcloud_exporter/issues/70


# Changelog for 1.2.0

The following sections list the changes for 1.2.0.

## Summary

 * Chg #53: Integrate standard web config
 * Chg #67: Add collector for server metrics

## Details

 * Change #53: Integrate standard web config

   We integrated the new web config from the Prometheus toolkit which provides a
   configuration for TLS support and also some basic builtin authentication. For
   the detailed configuration you check out the documentation.

   https://github.com/promhippie/hcloud_exporter/issues/53

 * Change #67: Add collector for server metrics

   Hetzner Cloud collects basic metrics on the hypervisor-level for each server. We
   have added a new collector which scrapes the latest available metric point for
   each running server. It is enabled by default.

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

   We have added a new optional collector, which is disabled by default, to gather
   metrics about the volumes part of the configured Hetzner Cloud project.

   https://github.com/promhippie/hcloud_exporter/issues/21

 * Change #24: Refactor build tools and project structure

   To have a unified project structure and build tooling we have integrated the
   same structure we already got within our GitHub exporter.

   https://github.com/promhippie/hcloud_exporter/issues/24

 * Change #25: Drop darwin/386 release builds

   We dropped the build of 386 builds on Darwin as this architecture is not
   supported by current Go versions anymore.

   https://github.com/promhippie/hcloud_exporter/issues/25

 * Change #39: Add collector for load balancers

   We have added a new optional collector, which is enabled by default, to gather
   metrics about all loadbalancers part of the configured Hetzner Cloud project.

   https://github.com/promhippie/hcloud_exporter/issues/39


# Changelog for 1.0.0

The following sections list the changes for 1.0.0.

## Summary

 * Chg #19: Add `vat` labels for net and gross values

## Details

 * Change #19: Add `vat` labels for net and gross values

   Added a new `vat` label for `gross` or `net` values to the `hcloud_server_price`
   metric. Depending on the setup this can be a breaking change and it may be
   necessary to adjust some dashboards and alerting rules.

   https://github.com/promhippie/hcloud_exporter/pull/19


# Changelog for 0.2.0

The following sections list the changes for 0.2.0.

## Summary

 * Chg #17: Add pricing collector
 * Chg #18: Add new metric to see if backups enabled

## Details

 * Change #17: Add pricing collector

   We added a new collector to gather information about the pricings, that way
   somebody could do calculations how much the costs are increasing or decreasing
   by sclae up or sclae down. The new collector includes new metrics named
   `hcloud_pricing_floating_ip`, `hcloud_pricing_image`,
   `hcloud_pricing_server_backup` and `hcloud_pricing_traffic`.

   https://github.com/promhippie/hcloud_exporter/pull/17

 * Change #18: Add new metric to see if backups enabled

   We added a new metric named `hcloud_server_backup` which indicates if a server
   got backups enabled or not, that way somebody could add some alerting if a
   server is missing a backup.

   https://github.com/promhippie/hcloud_exporter/pull/18


# Changelog for 0.1.1

The following sections list the changes for 0.1.1.

## Summary

 * Fix #11: Fix typo within `hcloud_server_incoming_traffic_bytes`
 * Chg #13: Pin go version to 1.10

## Details

 * Bugfix #11: Fix typo within `hcloud_server_incoming_traffic_bytes`

   We fixed a typo within the `hcloud_server_incoming_traffic_bytes` metric where
   we were just missing a tiny single letter.

   https://github.com/promhippie/hcloud_exporter/pull/11

 * Change #13: Pin go version to 1.10

   To make sure we got something nearly like reproducible builds and to fix the
   builds we should pin the build dependencies like the Go version to make sure it
   is always buildable.

   https://github.com/promhippie/hcloud_exporter/pull/13


# Changelog for 0.1.0

The following sections list the changes for 0.1.0.

## Summary

 * Chg #23: Initial release of basic version

## Details

 * Change #23: Initial release of basic version

   Just prepared an initial basic version which could be released to the public.

   https://github.com/promhippie/hcloud_exporter/issues/23


