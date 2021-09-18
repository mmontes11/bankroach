package cockroachdb

import (
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/golang-migrate/migrate/v4"
	migratecrdb "github.com/golang-migrate/migrate/v4/database/cockroachdb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const driverName = "postgres"

type crdb struct {
	*sql.DB
	*migrate.Migrate
}

func (db *crdb) MigrateUp() error {
	return runMigration(db.Up)
}

func (db *crdb) MigrateDown() error {
	return runMigration(db.Down)
}

func runMigration(migration func() error) error {
	err := migration()
	if err != migrate.ErrNoChange {
		return err
	}
	return nil
}

func New(dbUrl, migrationsFileUrl string) (*crdb, error) {
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
	migrate, err := migrate.NewWithDatabaseInstance(migrationsFileUrl, driverName, driver)
	if err != nil {
		return nil, err
	}

	return &crdb{DB: db, Migrate: migrate}, nil
}
