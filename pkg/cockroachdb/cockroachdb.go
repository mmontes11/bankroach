package cockroachdb

import (
	"context"
	"database/sql"

	"github.com/cockroachdb/cockroach-go/crdb"

	"github.com/golang-migrate/migrate/v4"
	migratecrdb "github.com/golang-migrate/migrate/v4/database/cockroachdb"

	// Postgres database driver: https://github.com/golang-migrate/migrate#databases
	_ "github.com/lib/pq"
	// Filesystem migration source driver: https://github.com/golang-migrate/migrate#migration-sources
	_ "github.com/golang-migrate/migrate/v4/source/file"
	// GitHub migration source driver: https://github.com/golang-migrate/migrate#migration-sources
	_ "github.com/golang-migrate/migrate/v4/source/github"
)

const driverName = "postgres"

type TxFn = func(*sql.Tx) error

type CRDB struct {
	*sql.DB
	*migrate.Migrate
}

func (db *CRDB) MigrateUp() error {
	return runMigration(db.Up)
}

func (db *CRDB) MigrateDown() error {
	return runMigration(db.Down)
}

func (db *CRDB) ExecuteTx(ctx context.Context, txOpts *sql.TxOptions, txFn TxFn) error {
	return crdb.ExecuteTx(ctx, db.DB, txOpts, txFn)
}

func runMigration(migration func() error) error {
	err := migration()
	if err != migrate.ErrNoChange {
		return err
	}
	return nil
}

func New(dbUrl, migrationsUrl string) (*CRDB, error) {
	db, err := sql.Open(driverName, dbUrl)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	driver, err := migratecrdb.WithInstance(db, &migratecrdb.Config{})
	if err != nil {
		return nil, err
	}
	migrate, err := migrate.NewWithDatabaseInstance(migrationsUrl, driverName, driver)
	if err != nil {
		return nil, err
	}

	return &CRDB{DB: db, Migrate: migrate}, nil
}
