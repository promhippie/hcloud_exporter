#!/bin/sh
set -e

if [ ! -d /var/lib/hcloud-exporter ]; then
    userdel hcloud-exporter 2>/dev/null || true
    groupdel hcloud-exporter 2>/dev/null || true
fi
