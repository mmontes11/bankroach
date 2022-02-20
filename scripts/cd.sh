#!/usr/bin/env bash

set -e

repo="mmontes"
release="bankroach"
chart="mmontes/bankroach"
namespace="bankroach"

helm repo add "$repo" https://charts.mmontes-dev.duckdns.org
helm repo update

echo "🚀 Deploying '$chart' with image version '$tag'..."
helm upgrade --install "$release" "$chart" --namespace "$namespace"