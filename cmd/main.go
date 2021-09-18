package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gotway/gotway/pkg/log"
	"github.com/mmontes11/bankroach/internal/controller"
	"github.com/mmontes11/bankroach/internal/repository"
	"github.com/mmontes11/bankroach/pkg/cockroachdb"
	"github.com/mmontes11/crypto-trade/pkg/scheduler"
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

	if err := db.MigrateDown(); err != nil {
		logger.Errorf("error rolling back migrations %v", err)
	}
	if err := db.MigrateUp(); err != nil {
		logger.Errorf("error executing migrations %v", err)
	}
	logger.Info("migrations completed")

	repo := repository.NewAccountRepo()
	ctrl := controller.New(repo, db, logger.WithField("type", "controller"))

	var wg sync.WaitGroup
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go func() {
			defer wg.Done()
			ctrl.CreateAccount(ctx, 1000)
		}()
	}
	wg.Wait()
	logger.Info("accounts created")

	s := scheduler.New(5*time.Second, func() {
		logger.Infof("starting %d transaction workers", 10)
		for i := 0; i < 10; i++ {
			go func() {
				if err := ctrl.Transfer(ctx); err != nil {
					logger.Errorf("error transferring %v", err)
				}
				if err := ctrl.PrintBalances(ctx); err != nil {
					logger.Errorf("error printing balances %v", err)
				}
			}()
		}
	})
	s.Start(ctx)
}
