Change: Read secrets form files

We have added proper support to load secrets like the password from files or
from base64-encoded strings. Just provide the flags or environment variables
for token or private key with a DSN formatted string like `file://path/to/file`
or `base64://Zm9vYmFy`.

https://github.com/promhippie/hcloud_exporter/pull/193
