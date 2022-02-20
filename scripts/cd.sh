#!/usr/bin/env bash

set -e

repo="bankroach"
release="bankroach"
chart="mmontes11/bankroach"
namespace="bankroach"

helm repo add "$repo" https://charts.mmontes-dev.duckdns.org
helm repo update

echo "🚀 Deploying '$chart' with image version '$tag'..."
helm upgrade --install "$release" "$chart" --namespace "$namespace"