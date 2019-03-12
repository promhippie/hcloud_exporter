# Changelog

## [1.0.0] - 2019-03-12

### Added

* Added a `vat` label (`gross` or `net`) to the `hcloud_server_price` metric. Depending on the setup this can be a breaking change and it may be necessary to adjust some dashboards and alerting rules.

## [0.2.0] - 2019-03-11

### Added

* Added new metric `hcloud_pricing_floating_ip`
* Added new metric `hcloud_pricing_image`
* Added new metric `hcloud_pricing_server_backup`
* Added new metric `hcloud_pricing_traffic`
* Added new metric `hcloud_server_backup`

### Fixed

* Fixed import path for golint command

## [0.1.1] - 2018-12-19

### Fixed

* Properly print floating IP values
* Pin Go to 1.10 to fix building
* Fix typo within `hcloud_server_incoming_traffic_bytes`

## [0.1.0] - 2018-10-06

### Added

* Initial release
