hcloud_floating_ip_active{<prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>}
: If 1 the floating IP is used by a server, 0 otherwise

hcloud_image_active{<prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>}
: If 1 the image is used by a server, 0 otherwise

hcloud_image_created_timestamp{<prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>}
: Timestamp when the image have been created

hcloud_image_deprecated_timestamp{<prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>}
: Timestamp when the image will be deprecated, 0 if not deprecated

hcloud_image_disk_bytes{<prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>}
: Size if the disk for the image in bytes

hcloud_image_size_bytes{<prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>}
: Size of the image in bytes

hcloud_pricing_floating_ip{<prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>}
: The cost of one floating IP per month

hcloud_pricing_image{<prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>}
: The cost of an image per GB/month

hcloud_pricing_loadbalancer_type{<prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>}
: The costs of a load balancer by type and location per month

hcloud_pricing_server_backup{}
: Will increase base server costs by specific percentage if server backups are enabled

hcloud_pricing_server_type{<prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>}
: The costs of a server by type and location per month

hcloud_pricing_traffic{<prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>}
: The cost of additional traffic per TB

hcloud_pricing_volume{<prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>}
: The cost of a volume per GB/month

hcloud_request_duration_seconds{collector}
: Histogram of latencies for requests to the api per collector

hcloud_request_failures_total{collector}
: Total number of failed requests to the api per collector

hcloud_server_backup{<prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>}
: If 1 server backups are enabled, 0 otherwise

hcloud_server_cores{<prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>}
: Server number of cores

hcloud_server_created_timestamp{<prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>}
: Timestamp when the server have been created

hcloud_server_disk_bytes{<prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>}
: Server disk in bytes

hcloud_server_included_traffic_bytes{<prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>}
: Included traffic for the server in bytes

hcloud_server_incoming_traffic_bytes{<prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>}
: Ingoing traffic to the server in bytes

hcloud_server_memory_bytes{<prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>}
: Server memory in bytes

hcloud_server_metrics_cpu{<prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>}
: Server CPU usage metric

hcloud_server_metrics_disk_read_bps{<prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>}
: Server disk write bytes/s metric

hcloud_server_metrics_disk_read_iops{<prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>}
: Server disk read iop/s metric

hcloud_server_metrics_disk_write_bps{<prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>}
: Server disk write bytes/s metric

hcloud_server_metrics_disk_write_iops{<prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>}
: Server disk write iop/s metric

hcloud_server_metrics_network_in_bps{<prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>}
: Server network incoming bytes/s metric

hcloud_server_metrics_network_in_pps{<prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>}
: Server network incoming packets/s metric

hcloud_server_metrics_network_out_bps{<prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>}
: Server network outgoing bytes/s metric

hcloud_server_metrics_network_out_pps{<prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>}
: Server network outgoing packets/s metric

hcloud_server_outgoing_traffic_bytes{<prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>}
: Outgoing traffic from the server in bytes

hcloud_server_price_hourly{<prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>}
: Price of the server billed hourly in €

hcloud_server_price_monthly{<prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>}
: Price of the server billed monthly in €

hcloud_server_running{<prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>}
: If 1 the server is running, 0 otherwise

hcloud_ssh_key{<prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>}
: Information about SSH keys in your HetznerCloud project

hcloud_volume_created{<prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>}
: Timestamp when the volume have been created

hcloud_volume_protection{<prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>}
: If 1 the volume is protected, 0 otherwise

hcloud_volume_size{<prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>}
: Size of the volume in GB

hcloud_volume_status{<prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>}
: If 1 the volume is availabel, 0 otherwise
