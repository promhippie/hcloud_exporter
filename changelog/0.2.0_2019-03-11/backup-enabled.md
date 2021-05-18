Change: Add new metric to see if backups enabled

We added a new metric named `hcloud_server_backup` which indicates if a server
got backups enabled or not, that way somebody could add some alerting if a
server is missing a backup.

https://github.com/promhippie/hcloud_exporter/pull/18
