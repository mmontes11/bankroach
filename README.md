# bankroach
Simple CRUD application using CockroachDB and Go

### Getting started
##### Run locally
```bash
docker compose up -d
```
```bash
make run
```
##### Kubernetes
```bash
helm repo add mmontes https://charts.mmontes-dev.duckdns.org
```
```bash
helm install cockroachdb-operator mmontes/cockroachdb-operator 
```
```bash
helm install bankroach mmontes/bankroach 
```