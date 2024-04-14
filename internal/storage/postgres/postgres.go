package postgres

import (
	"fmt"
	"github.com/fishmanDK/avito_test_task/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type BannerOperations interface {
	GetUserBanner() (struct{}, error)
}

type Postgres struct {
	db *sqlx.DB
}

func NewPostgres(cfg config.PostgresConfig) (*Postgres, error) {
	const op = "postgres.NewPostgres"
	db, err := sqlx.Open("postgres", cfg.String())
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	return &Postgres{
		db: db,
	}, nil
}
