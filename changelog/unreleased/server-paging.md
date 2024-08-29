Bugfix: Fetch metrics for all servers

For previous versions we have used the wrong client function to gather the list
of servers for the server metrics, this hvae been fixed by using a function that
automatically fetches all servers by iterating of the pagination.

https://github.com/promhippie/hcloud_exporter/issues/246
