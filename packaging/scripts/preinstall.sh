#!/bin/sh
set -e

if ! getent group hcloud-exporter >/dev/null 2>&1; then
    groupadd --system hcloud-exporter
fi

if ! getent passwd hcloud-exporter >/dev/null 2>&1; then
    useradd --system --create-home --home-dir /var/lib/hcloud-exporter --shell /bin/bash -g hcloud-exporter hcloud-exporter
fi
