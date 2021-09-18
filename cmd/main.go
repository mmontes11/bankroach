package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/gotway/gotway/pkg/log"
	"github.com/mmontes11/bankroach/pkg/cockroachdb"
)

func main() {
	dbUrl := flag.String(
		"db_url",
		"postgres://roach:@localhost:26257/bankroach?sslmode=disable",
		"CockroachDB URL connection string",
	)
	migrationsUrl := flag.String(
		"migrations_url",
		"file://migrations",
		"Migrations file connection string",
	)
	flag.Parse()

	logger := log.NewLogger(log.Fields{
		"service": "bankroach",
	}, "local", "debug", os.Stdout)

	ctx, cancel := signal.NotifyContext(context.Background(), []os.Signal{
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGKILL,
		syscall.SIGHUP,
		syscall.SIGQUIT}...,
	)
	defer cancel()

	db, err := cockroachdb.New(*dbUrl, *migrationsUrl)
	if err != nil {
		logger.Errorf("database error %v", err)
	}
	logger.Info("connected to database")

	logger.Info("running migrations")
	if err := db.MigrateUp(); err != nil {
		logger.Errorf("database migrations error %v", err)
	}
	logger.Info("migrations completed")

	<-ctx.Done()
}
