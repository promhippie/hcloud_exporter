Change: Add pricing collector

We added a new collector to gather information about the pricings, that way
somebody could do calculations how much the costs are increasing or decreasing
by sclae up or sclae down. The new collector includes new metrics named
`hcloud_pricing_floating_ip`, `hcloud_pricing_image`,
`hcloud_pricing_server_backup` and `hcloud_pricing_traffic`.

https://github.com/promhippie/hcloud_exporter/pull/17
