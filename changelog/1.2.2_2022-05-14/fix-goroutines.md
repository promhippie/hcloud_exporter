Bugfix: Another fix for go routines within server metrics

We disabled the server metrics by default for now until the implementation is
really stable to avoid any side effects. I have reintroduced routines, otherwise
the scrapetime will be far too high. This time I used wait groups to get
everything handled properly.

https://github.com/promhippie/hcloud_exporter/issues/74
