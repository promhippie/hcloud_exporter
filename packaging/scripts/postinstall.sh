#!/bin/sh
set -e

chown -R hcloud-exporter:hcloud-exporter /var/lib/hcloud-exporter
chmod 750 /var/lib/hcloud-exporter

if [ -d /run/systemd/system ]; then
    systemctl daemon-reload

    if systemctl is-enabled --quiet hcloud-exporter.service; then
        systemctl restart hcloud-exporter.service
    fi
fi
