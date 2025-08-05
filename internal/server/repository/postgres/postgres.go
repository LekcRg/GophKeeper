package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/LekcRg/GophKeeper/internal/config"
	"github.com/LekcRg/GophKeeper/internal/server/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type postgresRepo struct {
	db     *sqlx.DB
	config *config.Postgres
}

const PostgresDriver = "pgx"

func New(
	ctx context.Context, cfg *config.Postgres, log *zap.Logger,
) (*repository.Repository, error) {
	pg := &postgresRepo{
		config: cfg,
	}

	nativeDB, err := pg.getPgPool(ctx)
	if err != nil {
		return nil, err
	}

	db := sqlx.NewDb(nativeDB, PostgresDriver)

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Info("Success connect to postgres db")

	pg.db = db

	return &repository.Repository{
		DB:        pg,
		SQL:       nativeDB,
		UserRepo:  NewUserRepo(db),
		VaultRepo: NewVaultRepo(db),
	}, nil
}

func (p *postgresRepo) getURIFromConfig() string {
	dbURI := p.config.URI
	if p.config.URI == "" {
		dbURI = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?pool_max_conns=%s",
			p.config.User, p.config.Password, p.config.Host,
			p.config.Port, p.config.DB, p.config.MaxConns,
		)
	}

	return dbURI
}

func (p *postgresRepo) getPgPool(ctx context.Context) (*sql.DB, error) {
	connConfig, err := pgxpool.ParseConfig(p.getURIFromConfig())
	if err != nil {
		return nil, err
	}

	connPool, err := pgxpool.NewWithConfig(ctx, connConfig)
	if err != nil {
		return nil, err
	}

	return stdlib.OpenDBFromPool(connPool), nil
}

func (p *postgresRepo) Close() error {
	return p.db.Close()
}
