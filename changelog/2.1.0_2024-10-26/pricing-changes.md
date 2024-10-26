Change: Add type to IP pricing and add metrics for primary IPs

Since the client SDK has deprecated the previous handling for the pricing of IP
addresses we had to update the metrics to include the type and location of the
IPs. Besides that we have also added metrics for the pricing of the primary IP
addresses.

https://github.com/promhippie/hcloud_exporter/pull/272
