Change: Add `vat` labels for net and gross values

Added a new `vat` label for `gross` or `net` values to the `hcloud_server_price`
metric. Depending on the setup this can be a breaking change and it may be
necessary to adjust some dashboards and alerting rules.

https://github.com/promhippie/hcloud_exporter/pull/19
