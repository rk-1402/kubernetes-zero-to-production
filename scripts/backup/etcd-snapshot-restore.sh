#!/usr/bin/env bash
# Topic 64 / Capstone 10 — restore etcd from a snapshot (single-member recovery)
set -euo pipefail
systemctl stop kubelet
mv /var/lib/etcd /var/lib/etcd.old
ETCDCTL_API=3 etcdctl snapshot restore /tmp/snap.db \
  --name="${NODE_NAME:?set NODE_NAME}" \
  --initial-cluster="${NODE_NAME}=https://${NODE_IP:?set NODE_IP}:2380" \
  --initial-advertise-peer-urls="https://${NODE_IP}:2380" \
  --data-dir=/var/lib/etcd
systemctl start kubelet
