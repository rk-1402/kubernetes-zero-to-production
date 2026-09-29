#!/bin/sh
# Capstone 6 — logical Postgres backup to object storage (runs inside a CronJob container)
set -eu
TS=$(date +%Y%m%d-%H%M)
pg_dump -h pg-postgresql -U postgres shopdb | gzip > /tmp/backup-"$TS".sql.gz
aws s3 cp /tmp/backup-"$TS".sql.gz s3://acme-pg-backups/shopdb/backup-"$TS".sql.gz
