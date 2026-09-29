#!/usr/bin/env bash
# Topic 60 — use Kustomize as a Helm post-renderer:
#   helm install my-nginx bitnami/nginx --post-renderer ./post-renderer.sh
cat > all.yaml
kustomize build .
