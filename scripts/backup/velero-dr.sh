#!/usr/bin/env bash
# Capstone 10 — Velero cross-region backup & restore
set -euo pipefail
velero install \
  --provider aws --plugins velero/velero-plugin-for-aws:v1.9.0 \
  --bucket velero-dr --backup-location-config region=us-east-1 \
  --snapshot-location-config region=us-east-1 --secret-file ./credentials-velero
velero backup create shop-daily --include-namespaces=shop,db --wait
# In the DR cluster (us-west-2):
#   aws eks update-kubeconfig --name shop-dr --region us-west-2
#   velero restore create --from-backup shop-daily --wait
