package database

import (
	"database/sql"
	"fmt"
	"runtime"
	"strings"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"

	"bitbucket.org/junglee_games/getsetgo/pandora/log"
)

const dbDialect = "postgres"

type config struct {
	// MigrationTableName is the name of table used to track which version the database is in
	MigrationTableName string

	// Instrument if we want to send out data into elastic apm
	Instrument bool

	// Verbose adds a hook into the database that causes every query to be logged
	Verbose bool

	// Credentials used to connect to the database
	Credentials Credentials

	// ReadTimeout for socket reads.
	ReadTimeout time.Duration

	// WriteTimeout for socket writes
	WriteTimeout time.Duration

	DBOptions []bun.DBOption
}

var defaultConfig = config{
	MigrationTableName: "database_version",

	Instrument: true,
	Verbose:    false,

	Credentials: Credentials{},

	ReadTimeout:  time.Second * 30,
	WriteTimeout: time.Second * 30,
}

type OptFunc func(*config) error

func WithDisabledMigrations() OptFunc {
	return func(c *config) error {
		c.Credentials.MigrationsPath = ""

		return nil
	}
}

func WithDBName(name string) OptFunc {
	return func(c *config) error {
		c.Credentials.DBName = name

		return nil
	}
}

func WithInstrumentEnabled(enabled bool) OptFunc {
	return func(c *config) error {
		c.Instrument = enabled

		return nil
	}
}

func WithVerboseLogging(enabled bool) OptFunc {
	return func(c *config) error {
		c.Verbose = enabled

		return nil
	}
}

func WithDBOptions(options ...bun.DBOption) OptFunc {
	return func(c *config) error {
		c.DBOptions = options

		return nil
	}
}

type Credentials struct {
	Host           string
	Port           string
	User           string
	Pass           string
	DBName         string
	SchemaName     string
	MigrationsPath string
}

func defaultRdsCredentials() Credentials {
	return Credentials{
		Host: getWithDefault("POSTGRES_HOST", "localhost"),
		Port: getWithDefault("POSTGRES_PORT", "5432"),
		User: getWithDefault("POSTGRES_USER", "jungleegames"),
		Pass: getWithDefault("POSTGRES_PASSWORD", "jungleegames"),

		DBName:     getWithDefault("POSTGRES_DB", "jungleegames"),
		SchemaName: getWithDefault("POSTGRES_SCHEMA", "public"),

		MigrationsPath: findMigrationPath(),
	}
}

func defaultRedshiftCredentials() Credentials {
	return Credentials{
		Host: getWithDefault("REDSHIFT_HOST", "redshift.svc.prod.jungleegames.io"),
		Port: getWithDefault("REDSHIFT_PORT", "5439"),
		User: getWithDefault("REDSHIFT_USER", "jungleegames"),
		Pass: getWithDefault("REDSHIFT_PASSWORD", "jungleegames"),

		DBName:     getWithDefault("REDSHIFT_DB", "app_data_store"),
		SchemaName: getWithDefault("REDSHIFT_SCHEMA", "mobile_app_prod"),

		// We don't have migrations for redshift, we should but a lot of the data comes from segment.io
		// and we don't control that schema, in the future we would want to handle it some how.
		MigrationsPath: "",
	}
}

// Connect to a generic postgres server using the env variables
func Connect(opts ...OptFunc) (*bun.DB, error) {
	config := defaultConfig
	config.Credentials = defaultRdsCredentials()

	for _, opt := range opts {
		if err := opt(&config); err != nil {
			return nil, errors.Wrap(err, "failed to configure client")
		}
	}

	return connect(config)
}

func ConnectRedshift(opts ...OptFunc) (*bun.DB, error) {
	config := defaultConfig
	config.Credentials = defaultRedshiftCredentials()

	// override default timeout config for redshift
	config.ReadTimeout = 30 * time.Minute
	config.WriteTimeout = 30 * time.Minute

	for _, opt := range opts {
		if err := opt(&config); err != nil {
			return nil, errors.Wrap(err, "failed to configure client")
		}
	}

	return connect(config)
}

func connect(config config) (*bun.DB, error) {
	addr := config.Credentials.Host
	if !strings.Contains(config.Credentials.Host, ":") {
		addr = fmt.Sprintf("%s:%s", config.Credentials.Host, config.Credentials.Port)
	}

	driver := pgdriver.NewConnector(
		pgdriver.WithNetwork("tcp"),
		pgdriver.WithAddr(addr),
		pgdriver.WithUser(config.Credentials.User),
		pgdriver.WithPassword(config.Credentials.Pass),
		pgdriver.WithDatabase(config.Credentials.DBName),
		pgdriver.WithInsecure(true),
		pgdriver.WithConnParams(map[string]interface{}{
			"search_path": config.Credentials.SchemaName,
		}),
		pgdriver.WithReadTimeout(config.ReadTimeout),
		pgdriver.WithWriteTimeout(config.WriteTimeout),
	)

	conn := sql.OpenDB(driver)
	conn.SetMaxOpenConns(runtime.NumCPU() * 3)
	conn.SetConnMaxIdleTime(time.Second * 60 * 3)

	db := bun.NewDB(conn, pgdialect.New(), config.DBOptions...)

	// We do this to ensure the connection works
	if _, err := db.Exec("SELECT 1"); err != nil {
		return nil, err
	}

	if config.Instrument {
		db.AddQueryHook(&NewrelicHook{
			Host:     config.Credentials.Host,
			Port:     config.Credentials.Port,
			Database: config.Credentials.DBName,
		})
	}

	if config.Verbose {
		db.AddQueryHook(verboseLogger{})
	}

	if config.Credentials.MigrationsPath != "" {
		if _, _, _, err := RunSchemaMigrations(config.MigrationTableName, config.Credentials); err != nil {
			return db, errors.Wrap(err, "failed to migrate")
		}
	} else {
		log.
			WithField("host", config.Credentials.Host).
			Infof("No migrations paths specified will not run migrations")
	}

	return db, nil
}
