#!/usr/bin/env bash
# Capstone 3 — alternative HTTP load generator
# Port-forward first: kubectl -n shop port-forward svc/backend 8080:8080 &
hey -z 5m -c 200 -q 50 http://localhost:8080/api/items
