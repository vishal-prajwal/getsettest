package database

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
)

func truncateAllTables(t *testing.T, db *bun.DB, migrationTableName string) {
	require.NotNil(t, db, "*bun.DB is nil, a proper instance is required")

	var truncateQueries []string

	// Purge all the tables with the cascading option
	rows, err := db.Query(fmt.Sprintf("SELECT 'TRUNCATE ' || tablename || ' CASCADE;' FROM pg_tables WHERE tablename NOT LIKE 'pg_%%' AND tablename NOT LIKE 'sql_%%' AND tablename NOT LIKE '%s';", migrationTableName))
	if err != nil {
		require.NoError(t, err, "failed to list all the existing tables")
	}

	err = db.ScanRows(context.TODO(), rows, &truncateQueries)
	if err != nil {
		require.NoError(t, err, "failed to list all the existing tables")
	}

	truncateQuery := strings.Join(truncateQueries, "\n")

	// nothing to do here
	if truncateQuery == "" {
		return
	}

	_, err = db.Exec(truncateQuery)
	require.NoError(t, err, fmt.Sprintf("failed to execute: %s", truncateQuery))
}

func CreateNewTestDatabase(t *testing.T, name string, opts ...OptFunc) *bun.DB {
	config := defaultConfig

	for _, opt := range opts {
		require.NoError(t, opt(&config))
	}

	db, err := Connect(append(opts, WithDisabledMigrations())...)
	require.NoError(t, err)

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s;", name))
	if err != nil {
		require.Error(t, err, `ERROR #42P04 database "%s" already exists`, name)
	}

	require.NoError(t, db.Close())

	// We need to reconnect because bun doesn't support changing the database after connecting
	db, err = Connect(append(opts, WithDBName(name))...)
	require.NoError(t, err)

	truncateAllTables(t, db, config.MigrationTableName)

	return db
}
