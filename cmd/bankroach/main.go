package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gotway/gotway/pkg/log"
	"github.com/mmontes11/bankroach/internal/config"
	"github.com/mmontes11/bankroach/internal/controller"
	"github.com/mmontes11/bankroach/internal/repository"
	"github.com/mmontes11/bankroach/pkg/cockroachdb"
	"github.com/mmontes11/crypto-trade/pkg/scheduler"
)

func main() {
	conf := config.Get()
	logger := log.NewLogger(log.Fields{
		"service": "bankroach",
	}, conf.Env, conf.LogLevel, os.Stdout)

	ctx, cancel := signal.NotifyContext(context.Background(), []os.Signal{
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGKILL,
		syscall.SIGHUP,
		syscall.SIGQUIT}...,
	)
	defer cancel()

	db, err := cockroachdb.New(conf.DBUrl, conf.MigrationsUrl)
	if err != nil {
		logger.Fatalf("database error %v", err)
	}
	logger.Info("connected to database")

	if err := db.MigrateDown(); err != nil {
		logger.Fatalf("error rolling back migrations %v", err)
	}
	if err := db.MigrateUp(); err != nil {
		logger.Fatalf("error executing migrations %v", err)
	}
	logger.Info("migrations completed")

	repo := repository.NewAccountRepo()
	ctrl := controller.New(repo, db, conf, logger.WithField("type", "controller"))

	var wg sync.WaitGroup
	wg.Add(conf.NumAccounts)
	for i := 0; i < conf.NumAccounts; i++ {
		go func() {
			defer wg.Done()
			ctrl.CreateAccount(ctx, int64(conf.InitialBalance))
		}()
	}
	wg.Wait()
	logger.Info("accounts created")

	s := scheduler.New(conf.ScheduleInterval, func() {
		logger.Infof("starting %d transaction workers", conf.NumWorkers)
		for i := 0; i < conf.NumWorkers; i++ {
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
