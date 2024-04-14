package service

import (
	"fmt"
	"github.com/fishmanDK/avito_test_task/internal/storage"
	"github.com/fishmanDK/avito_test_task/models"
)

type BannerManager struct {
	storage *storage.Storage
}

func NewBannerManager(storage *storage.Storage) *BannerManager {
	return &BannerManager{
		storage: storage,
	}
}

func (bm *BannerManager) GetBanners(params models.GetAllBannersParams) (*[]models.BannerWithDetails, error) {
	const op = "service.GetAllBanners"

	banners, err := bm.storage.DB.GetBannersWithDetails(params)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	return &banners, nil
}

func (bm *BannerManager) CreateBanner(newBanner models.CreateBannerRequest) error {
	const op = "service.CreateBanner"

	err := bm.storage.DB.CreateBanner(newBanner)
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	return nil
}

func (bm *BannerManager) ChangeBanner(bannerID int64, req models.ChangeBannerRequest) error {
	const op = "service.ChangeBanner"

	err := bm.storage.DB.CheckBanner(bannerID, 0, 0)
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	err = bm.storage.DB.ChangeBanner(bannerID, req)
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	return nil
}

//func (bm *BannerManager) deleteFullBanner(bannerID int64) error {
//	const op = "service.DeleteBanner"
//
//	err := bm.storage.DB.DeleteBanner(bannerID)
//	if err != nil {
//		return fmt.Errorf("%s: %v", op, err)
//	}
//	return nil
//}
