hcloud_floating_ip_active{id, server, location, type, ip}
: If 1 the floating IP is used by a server, 0 otherwise

hcloud_image_active{id, name, type, server, flavor, version}
: If 1 the image is used by a server, 0 otherwise

hcloud_image_created_timestamp{id, name, type, server, flavor, version}
: Timestamp when the image have been created

hcloud_image_deprecated_timestamp{id, name, type, server, flavor, version}
: Timestamp when the image will be deprecated, 0 if not deprecated

hcloud_image_disk_bytes{id, name, type, server, flavor, version}
: Size if the disk for the image in bytes

hcloud_image_size_bytes{id, name, type, server, flavor, version}
: Size of the image in bytes

hcloud_pricing_floating_ip{currency, vat, type, location}
: The cost of one floating IP per month

hcloud_pricing_image{currency, vat}
: The cost of an image per GB/month

hcloud_pricing_loadbalancer_type{currency, vat, type, location}
: The costs of a load balancer by type and location per month

hcloud_pricing_loadbalancer_type_traffic{currency, vat, type, location}
: The costs of additional traffic per TB for a load balancer by type and location per month

hcloud_pricing_primary_ip{currency, vat, type, location}
: The cost of one primary IP per month

hcloud_pricing_server_backup{}
: Will increase base server costs by specific percentage if server backups are enabled

hcloud_pricing_server_type{currency, vat, type, location}
: The costs of a server by type and location per month

hcloud_pricing_server_type_traffic{currency, vat, type, location}
: The costs of additional traffic per TB for a server by type and location per month

hcloud_pricing_volume{currency, vat}
: The cost of a volume per GB/month

hcloud_request_duration_seconds{collector}
: Histogram of latencies for requests to the api per collector

hcloud_request_failures_total{collector}
: Total number of failed requests to the api per collector

hcloud_server_backup{id, name, datacenter}
: If 1 server backups are enabled, 0 otherwise

hcloud_server_cores{id, name, datacenter}
: Server number of cores

hcloud_server_created_timestamp{id, name, datacenter}
: Timestamp when the server have been created

hcloud_server_disk_bytes{id, name, datacenter}
: Server disk in bytes

hcloud_server_included_traffic_bytes{id, name, datacenter}
: Included traffic for the server in bytes

hcloud_server_incoming_traffic_bytes{id, name, datacenter}
: Ingoing traffic to the server in bytes

hcloud_server_memory_bytes{id, name, datacenter}
: Server memory in bytes

hcloud_server_metrics_cpu{id, name, datacenter}
: Server CPU usage metric

hcloud_server_metrics_disk_read_bps{id, name, datacenter, disk}
: Server disk write bytes/s metric

hcloud_server_metrics_disk_read_iops{id, name, datacenter, disk}
: Server disk read iop/s metric

hcloud_server_metrics_disk_write_bps{id, name, datacenter, disk}
: Server disk write bytes/s metric

hcloud_server_metrics_disk_write_iops{id, name, datacenter, disk}
: Server disk write iop/s metric

hcloud_server_metrics_network_in_bps{id, name, datacenter, interface}
: Server network incoming bytes/s metric

hcloud_server_metrics_network_in_pps{id, name, datacenter, interface}
: Server network incoming packets/s metric

hcloud_server_metrics_network_out_bps{id, name, datacenter, interface}
: Server network outgoing bytes/s metric

hcloud_server_metrics_network_out_pps{id, name, datacenter, interface}
: Server network outgoing packets/s metric

hcloud_server_outgoing_traffic_bytes{id, name, datacenter}
: Outgoing traffic from the server in bytes

hcloud_server_price_hourly{id, name, datacenter, vat}
: Price of the server billed hourly in €

hcloud_server_price_monthly{id, name, datacenter, vat}
: Price of the server billed monthly in €

hcloud_server_running{id, name, datacenter}
: If 1 the server is running, 0 otherwise

hcloud_ssh_key{id, name, fingerprint}
: Information about SSH keys in your Hetzner Cloud project

hcloud_volume_created{id, server, location, name}
: Timestamp when the volume have been created

hcloud_volume_protection{id, server, location, name}
: If 1 the volume is protected, 0 otherwise

hcloud_volume_size{id, server, location, name}
: Size of the volume in GB

hcloud_volume_status{id, server, location, name}
: If 1 the volume is availabel, 0 otherwise
