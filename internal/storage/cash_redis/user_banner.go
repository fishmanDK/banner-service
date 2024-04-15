package cash_redis

import (
	"encoding/json"
	"fmt"
	"github.com/fishmanDK/avito_test_task/models"
	"log"
	"time"
)

const (
	bannerCashExpiration = 5 * time.Minute
)

func (r *CashRedis) CashGetUserBanner() (*models.BannerWithDetails, error) {
	var cashBanner *models.BannerWithDetails
	return cashBanner, nil
}

func (r *CashRedis) SaveBanner(banners ...*models.BannerWithDetails) error {
	const op = "cash_redis.SaveBanner"
	bannerJSON, err := json.Marshal(banners)
	if err != nil {
		return err
	}
	for i := range banners {
		err = r.client.Set(fmt.Sprintf("bannerID:%d", banners[i].BannerID), bannerJSON, bannerCashExpiration).Err()
		if err != nil {
			log.Printf("%s: failed cash banner id = %d", op, banners[i].BannerID)
		}
	}

	return nil
}
