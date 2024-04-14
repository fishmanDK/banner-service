package service

import (
	"context"
	"fmt"
	trigger_hook_grpc "github.com/fishmanDK/avito_test_task/internal/clients/trigger_service/grpc"
	"github.com/fishmanDK/avito_test_task/internal/config"
	"github.com/fishmanDK/avito_test_task/internal/storage"
	"github.com/fishmanDK/avito_test_task/models"
	"log/slog"
)

type Service struct {
	Auth
	UserBannerGetter
	BannerOperations
	ParamsOperations
	DeleteService *DeleteService
}

type ParamsOperations interface {
	CreateTag(tag models.Tag) error
	CreateFeature(feature models.Feature) error
}

type Auth interface {
	Authentication(user models.User) (models.Token, error)
	CreateUser(newUser models.NewUser) error
	ParseToken(accessToken string) (string, error)
}

type UserBannerGetter interface {
	GetUserBanner(params models.UserBanner) (*models.BannerWithDetails, error)
}

type BannerOperations interface {
	ChangeBanner(bannerID int64, req models.ChangeBannerRequest) error
	GetBanners(req models.GetAllBannersParams) (*[]models.BannerWithDetails, error)
	CreateBanner(req models.CreateBannerRequest) error
}

func NewService(logger *slog.Logger, cfg config.Clients, storage *storage.Storage) (*Service, error) {
	const op = "service.NewService"
	triggerHookClient, err := trigger_hook_grpc.NewClient(
		context.Background(),
		logger,
		cfg.TriggerHookService.Address,
		cfg.TriggerHookService.Timeout,
		cfg.TriggerHookService.RetriesCount)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	deleteService := NewDeleteService(triggerHookClient, storage)

	return &Service{
		Auth:             NewAuthService(storage),
		UserBannerGetter: NewUserBannerManager(storage),
		BannerOperations: NewBannerManager(storage),
		ParamsOperations: NewParamsManager(storage),
		DeleteService:    deleteService,
	}, nil
}
