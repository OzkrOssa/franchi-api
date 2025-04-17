package postgres

import (
	"context"
	"fmt"

	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/OzkrOssa/franchi-api/adapter/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

/**
 * DB is a wrapper for PostgreSQL database connection
 * that uses pgxpool as database driver.
 * It also holds a reference to squirrel.StatementBuilderType
 * which is used to build SQL queries that compatible with PostgreSQL syntax
 */
type DB struct {
	*pgxpool.Pool
	QueryBuilder *squirrel.StatementBuilderType
}

// New creates a new PostgreSQL database instance
func New(ctx context.Context, config *config.DB) (*DB, error) {
	url := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable",
		config.Connection,
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Name,
	)

	db, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, err
	}

	err = db.Ping(ctx)
	if err != nil {
		return nil, err
	}

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	return &DB{
		db,
		&psql,
	}, nil
}

func (db *DB) Migrate(ctx context.Context) error {
	sqlDB, err := sql.Open("pgx", db.Pool.Config().ConnString())
	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://adapter/storage/postgres/migrations",
		"postgres", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func (db *DB) ErrorCode(err error) string {
	pgErr := err.(*pgconn.PgError)
	return pgErr.Code
}

func (db *DB) Close() {
	db.Pool.Close()
}
