Bugfix: Correctly read loadbalancer traffic

We used a wrong attribute to read the loadbalancer traffic which resulted in
missing metrics for the realtime traffic in and out for all loadbalancers. With
this fix you should be able to use the metrics.

https://github.com/promhippie/hcloud_exporter/issues/175
