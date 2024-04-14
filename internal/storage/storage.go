package storage

import (
	"fmt"
	"github.com/fishmanDK/avito_test_task/internal/config"
	"github.com/fishmanDK/avito_test_task/internal/storage/cash_redis"
	"github.com/fishmanDK/avito_test_task/internal/storage/postgres"
)

type Storage struct {
	Cash *cash_redis.CashRedis
	DB   *postgres.Postgres
}

func MustStorage(cfg config.PostgresConfig) (*Storage, error) {
	const op = "storage.MustStorage"

	redis_cash, err := cash_redis.NewCashRedis()
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	postgre, err := postgres.NewPostgres(cfg)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	return &Storage{
		Cash: redis_cash,
		DB:   postgre,
	}, nil
}
