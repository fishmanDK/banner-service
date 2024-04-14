package service

import (
	"context"
	"fmt"
	trigger_hook_grpc "github.com/fishmanDK/avito_test_task/internal/clients/trigger_service/grpc"
	"github.com/fishmanDK/avito_test_task/internal/storage"
)

type DeleteService struct {
	storage           *storage.Storage
	triggerHookClient *trigger_hook_grpc.Client
}

func NewDeleteService(triggerHookClient *trigger_hook_grpc.Client, storage *storage.Storage) *DeleteService {
	return &DeleteService{
		triggerHookClient: triggerHookClient,
		storage:           storage,
	}
}

func (d *DeleteService) ScheduleDeleteBanner(ctx context.Context, bannerID, tagID, featureID int64) error {
	const op = "service.ScheduleDeleteBanner"

	err := d.storage.DB.CheckBanner(bannerID, tagID, featureID)
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	status, err := d.triggerHookClient.ScheduleDeletion(ctx, bannerID, tagID, featureID)
	if err != nil || !status {
		return fmt.Errorf("%s: %v", op, err)
	}

	return nil
}

func (bm *DeleteService) DeleteBanner(bannerID int64) error {
	const op = "service.DeleteBanner"

	err := bm.storage.DB.DeleteBanner(bannerID)
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	return nil
}

func (bm *DeleteService) DeleteBannerByParams(tagID, featuresID int64) error {
	const op = "service.DeleteBanner"

	err := bm.storage.DB.DeleteBannerByParams(tagID, featuresID)
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	return nil
}
