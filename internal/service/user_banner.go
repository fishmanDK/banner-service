package service

import (
	"context"
	"fmt"
	"github.com/fishmanDK/avito_test_task/internal/storage"
	"github.com/fishmanDK/avito_test_task/internal/storage/cash_redis"
	"github.com/fishmanDK/avito_test_task/models"
)

type UserBannerManager struct {
	storage *storage.Storage
	cash    *cash_redis.CashRedis
}

func NewUserBannerManager(storage *storage.Storage, cash *cash_redis.CashRedis) *UserBannerManager {
	return &UserBannerManager{
		storage: storage,
		cash:    cash,
	}
}

func (ubm *UserBannerManager) GetUserBanner(ctx context.Context, params models.UserBanner) (*models.BannerWithDetails, error) {
	const op = "service.GetUserBanner"
	if params.UseLastRevision {
		banner, err := ubm.storage.DB.GetUserBanner(params)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		err = ubm.cash.SaveBanner(banner)
		if err != nil {
			fmt.Errorf("%s: %w", op, err)
		}

		return banner, nil
	} else {
		cashBanner, err := ubm.storage.Cash.CashGetUserBanner()
		if err != nil {
			banner, err := ubm.storage.DB.GetUserBanner(params)
			if err != nil {
				return nil, fmt.Errorf("%s: %v", op, err)
			}

			err = ubm.cash.SaveBanner(banner)
			if err != nil {
				fmt.Errorf("%s: %v", op, err)
			}
			return banner, err
		}
		return cashBanner, err
	}
}
