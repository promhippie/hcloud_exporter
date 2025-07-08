#!/bin/sh
set -e

systemctl stop hcloud-exporter.service || true
systemctl disable hcloud-exporter.service || true
