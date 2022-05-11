Change: Add collector for server metrics

Hetzner Cloud collects basic metrics on the hypervisor-level for each server. We
have added a new collector which scrapes the latest available metric point for
each running server. It is enabled by default.

https://github.com/promhippie/hcloud_exporter/pull/67
