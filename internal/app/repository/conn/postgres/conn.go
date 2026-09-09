package rcpostgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/migrate"

	"github.com/mahkamovvlad/catalog-service/internal/app/config/section"
	"github.com/mahkamovvlad/catalog-service/migration"
)

type (
	Client struct {
		_bunDB
		rawBunDB *bun.DB

		cfg section.RepositoryPostgres
	}

	_bunDB = bun.IDB
)

func (c *Client) GetRawBunDB() *bun.DB {
	return c.rawBunDB
}

func NewClient(ctx context.Context, cfg section.RepositoryPostgres) (*Client, error) {
	var u url.URL
	u.Scheme = "postgres"
	u.Host = cfg.Address
	u.User = url.UserPassword(cfg.Username, cfg.Password)
	u.Path = cfg.Name

	q := u.Query()
	q.Set("sslmode", "disable")
	u.RawQuery = q.Encode()

	dsn := u.String()

	log.Printf("Connecting to PostgreSQL at %s (Database: %s)...", cfg.Address, cfg.Name)

	connector := pgdriver.NewConnector(
		pgdriver.WithDSN(dsn),
		pgdriver.WithTimeout(cfg.ConnTimeout),
		pgdriver.WithReadTimeout(cfg.ReadTimeout),
		pgdriver.WithWriteTimeout(cfg.WriteTimeout),
	)

	sqlDB := sql.OpenDB(connector)
	sqlDB.SetMaxOpenConns(10)

	bunDB := bun.NewDB(sqlDB, pgdialect.New(), bun.WithDiscardUnknownColumns())

	ctxPing, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := bunDB.PingContext(ctxPing); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("PostgreSQL connection verified successfully via Ping")

	return &Client{
		_bunDB:   bunDB,
		rawBunDB: bunDB,
		cfg:      cfg,
	}, nil
}

func (c *Client) Migrate(ctx context.Context) (oldVer, newVer int64, err error) {
	migrations := migrate.NewMigrations()
	if err := migrations.Discover(migration.Postgres); err != nil {
		return 0, 0, fmt.Errorf("failed to discover migrations: %w", err)
	}

	migrator := migrate.NewMigrator(
		c.rawBunDB,
		migrations,
		migrate.WithTableName(c.cfg.MigrationTable),
		migrate.WithLocksTableName(c.cfg.MigrationTable+"_lock"),
		migrate.WithMarkAppliedOnSuccess(true),
	)
	if err := migrator.Init(ctx); err != nil {
		return 0, 0, fmt.Errorf("failed to init migrator: %w", err)
	}

	applied, err := migrator.AppliedMigrations(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to load applied migrations: %w", err)
	}

	for _, m := range applied {
		v, err := strconv.ParseInt(m.Name, 10, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("failed to parse applied migration version %q: %w", m.Name, err)
		}
		if v > oldVer {
			oldVer = v
		}
	}

	group, err := migrator.Migrate(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to execute migrations: %w", err)
	}

	newVer = oldVer
	for _, m := range group.Migrations {
		v, err := strconv.ParseInt(m.Name, 10, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("failed to parse new migration version %q: %w", m.Name, err)
		}
		if v > newVer {
			newVer = v
		}
	}

	return oldVer, newVer, nil
}
