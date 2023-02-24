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
helm repo add mmontes https://mmontes11.github.io/charts
helm install cockroachdb-operator mmontes/cockroachdb-operator 
helm install bankroach mmontes/bankroach 
```
