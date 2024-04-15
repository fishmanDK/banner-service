package service

import (
	"context"
	"fmt"
	"github.com/fishmanDK/avito_test_task/internal/storage"
	"github.com/fishmanDK/avito_test_task/internal/storage/cash_redis"
	"github.com/fishmanDK/avito_test_task/models"
)

type BannerManager struct {
	storage *storage.Storage
	cash    *cash_redis.CashRedis
}

func NewBannerManager(storage *storage.Storage, cash *cash_redis.CashRedis) *BannerManager {
	return &BannerManager{
		storage: storage,
		cash:    cash,
	}
}

func (bm *BannerManager) GetBanners(ctx context.Context, params models.GetAllBannersParams) ([]*models.BannerWithDetails, error) {
	const op = "service.GetAllBanners"

	banners, err := bm.storage.DB.GetBannersWithDetails(params)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = bm.cash.SaveBanner(banners...)
	return banners, nil
}

func (bm *BannerManager) CreateBanner(ctx context.Context, newBanner models.CreateBannerRequest) error {
	const op = "service.CreateBanner"

	err := bm.storage.DB.CreateBanner(newBanner)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (bm *BannerManager) ChangeBanner(ctx context.Context, bannerID int64, req models.ChangeBannerRequest) error {
	const op = "service.ChangeBanner"

	err := bm.storage.DB.CheckBanner(bannerID, 0, 0)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = bm.storage.DB.ChangeBanner(bannerID, req)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
