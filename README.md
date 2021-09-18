# bankroach
Simple CRUD application using CockroachDB and Go

```bash
go run cmd/main.go -db_url='postgres://roach:@localhost:26257/bankroach?sslmode=disable' -migrations_url='file://migrations'
```