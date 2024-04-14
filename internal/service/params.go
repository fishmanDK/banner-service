package service

import (
	"fmt"
	"github.com/fishmanDK/avito_test_task/internal/storage"
	"github.com/fishmanDK/avito_test_task/models"
)

type ParamsManager struct {
	storage *storage.Storage
}

func NewParamsManager(storage *storage.Storage) *ParamsManager {
	return &ParamsManager{
		storage: storage,
	}
}

func (pm *ParamsManager) CreateTag(tag models.Tag) error {
	const op = "service.CreateTag"

	err := pm.storage.DB.CreateTag(tag)
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	return nil
}

func (pm *ParamsManager) CreateFeature(feature models.Feature) error {
	const op = "service.CreateFeature"

	err := pm.storage.DB.CreateFeature(feature)
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	return nil
}
