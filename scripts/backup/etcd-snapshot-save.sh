#!/usr/bin/env bash
# Topic 64 / Capstone 10 — take an etcd snapshot and copy off-cluster (kubeadm control-plane node)
set -euo pipefail
ETCDCTL_API=3 etcdctl snapshot save /tmp/snap.db \
  --endpoints=https://127.0.0.1:2379 \
  --cacert=/etc/kubernetes/pki/etcd/ca.crt \
  --cert=/etc/kubernetes/pki/etcd/server.crt \
  --key=/etc/kubernetes/pki/etcd/server.key
aws s3 cp /tmp/snap.db "s3://etcd-dr/snap-$(date +%F-%H%M).db"
