package service

import (
	"context"
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

func (pm *ParamsManager) CreateTag(ctx context.Context, tag models.Tag) error {
	const op = "service.CreateTag"

	err := pm.storage.DB.CreateTag(tag)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (pm *ParamsManager) CreateFeature(ctx context.Context, feature models.Feature) error {
	const op = "service.CreateFeature"

	err := pm.storage.DB.CreateFeature(feature)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
