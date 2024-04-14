package service

import (
	"github.com/fishmanDK/avito_test_task/internal/storage"
	"github.com/fishmanDK/avito_test_task/models"
)

type UserBannerManager struct {
	storage *storage.Storage
}

func NewUserBannerManager(storage *storage.Storage) *UserBannerManager {
	return &UserBannerManager{
		storage: storage,
	}
}

func (ubm *UserBannerManager) GetUserBanner(params models.UserBanner) (*models.BannerWithDetails, error) {
	if params.UseLastRevision {
		banner, err := ubm.storage.DB.GetUserBanner(params)
		if err != nil {
			// TODO: logger
			return nil, err
		}

		return banner, nil
	} else {
		cash_banner, err := ubm.storage.Cash.CashGetUserBanner()
		if err != nil {
			// TODO: logger
			// если не нашел в кеше, то идем в базу
			banner, err := ubm.storage.DB.GetUserBanner(params)
			if err != nil {
				// TODO: logger
				return nil, err
			}
			return banner, err
		}
		// TODO: если нет в кеше то идем в базу(сделать это правильно)
		return cash_banner, err

	}
}
