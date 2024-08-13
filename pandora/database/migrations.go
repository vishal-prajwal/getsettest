package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/pkg/errors"
)

func RunSchemaMigrations(migrationVersionTable string, c Credentials) (oldVersion uint, newVersion uint, dirty bool, err error) {
	url := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable", dbDialect, c.User, c.Pass, c.Host, c.Port, c.DBName)

	dbD, err := sql.Open(dbDialect, url)
	if err != nil {
		return 0, 0, false, errors.Wrapf(err, "failed to connect to data store: %q", url)
	}

	driver, err := postgres.WithInstance(dbD, &postgres.Config{MigrationsTable: migrationVersionTable})
	if err != nil {
		return 0, 0, false, errors.Wrap(err, "failed to initialize driver")
	}

	m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s", c.MigrationsPath), "", driver)
	if err != nil {
		return 0, 0, false, err
	}

	oldVersion, _, _ = m.Version()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return 0, 0, false, errors.WithStack(err)
	}

	newVersion, dirty, _ = m.Version()

	_ = driver.Unlock()
	_ = driver.Close()

	return
}

func findMigrationPath() string {
	migrations := os.Getenv("POSTGRES_MIGRATIONS_PATH")
	if migrations != "" {
		return migrations
	}

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	for {
		if dir == "/" || dir == "." || dir == "gitlab.com" {
			return ""
		}

		_, err := os.Stat(dir + "/migrations")
		if err == nil {
			return dir + "/migrations"
		}

		dir = filepath.Dir(dir)
	}
}
