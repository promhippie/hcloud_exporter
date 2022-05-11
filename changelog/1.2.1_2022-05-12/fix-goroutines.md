Bugfix: Fix go routine errors within server metrics

We fixed a go routines issue within the new server metrics. We just got rid of
the routines to avoid any errors related to sending to closed channels.

https://github.com/promhippie/hcloud_exporter/issues/70
