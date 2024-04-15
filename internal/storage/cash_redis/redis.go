package cash_redis

import (
	"fmt"
	"github.com/fishmanDK/avito_test_task/models"
	"github.com/go-redis/redis"
)

type CashBannerOperations interface {
	CashGetUserBanner() (*models.BannerWithDetails, error)
}

type CashRedis struct {
	client *redis.Client
}

func NewCashRedis() (*CashRedis, error) {
	// TODO: config redis
	const op = "cash_redis.NewCashRedis"
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	return &CashRedis{
		client: client,
	}, nil
}
