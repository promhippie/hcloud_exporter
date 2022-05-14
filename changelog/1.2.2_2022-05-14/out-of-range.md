Bugfix: Fix index out of range issue within server metrics

The code has not checked if an index have been really available within the
server metrics API response. With this fix it gets properly handled.

https://github.com/promhippie/hcloud_exporter/issues/72
