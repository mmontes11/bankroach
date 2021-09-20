package config

import (
	"time"

	"github.com/gotway/gotway/pkg/env"
)

type Config struct {
	DBUrl         string
	MigrationsUrl string

	NumAccounts    int
	InitialBalance int

	ScheduleInterval time.Duration
	NumWorkers       int

	LogLevel string
	Env      string
}

func Get() Config {
	return Config{
		DBUrl: env.Get(
			"DB_URL",
			"postgres://roach:@localhost:26257/bankroach?sslmode=disable",
		),
		MigrationsUrl: env.Get("MIGRATIONS_URL", "file://migrations"),

		NumAccounts:    env.GetInt("NUM_ACCOUNTS", 5),
		InitialBalance: env.GetInt("INITIAL_BALANCE", 1000),

		ScheduleInterval: env.GetDuration("SCHEDULE_INTERVAL_SECONDS", 5) * time.Second,
		NumWorkers:       env.GetInt("NUM_WORKERS", 10),

		Env:      env.Get("ENV", "local"),
		LogLevel: env.Get("LOG_LEVEL", "debug"),
	}
}
