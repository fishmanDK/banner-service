package service

import (
	"context"
	"fmt"
	trigger_hook_grpc "github.com/fishmanDK/avito_test_task/internal/clients/trigger_service/grpc"
	"github.com/fishmanDK/avito_test_task/internal/config"
	"github.com/fishmanDK/avito_test_task/internal/storage"
	"github.com/fishmanDK/avito_test_task/internal/storage/cash_redis"
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
	CreateTag(ctx context.Context, tag models.Tag) error
	CreateFeature(ctx context.Context, feature models.Feature) error
}

type Auth interface {
	Authentication(ctx context.Context, user models.User) (models.Token, error)
	CreateUser(ctx context.Context, newUser models.NewUser) error
	ParseToken(accessToken string) (string, error)
}

type UserBannerGetter interface {
	GetUserBanner(ctx context.Context, params models.UserBanner) (*models.BannerWithDetails, error)
}

type BannerOperations interface {
	ChangeBanner(ctx context.Context, bannerID int64, req models.ChangeBannerRequest) error
	GetBanners(ctx context.Context, req models.GetAllBannersParams) ([]*models.BannerWithDetails, error)
	CreateBanner(ctx context.Context, req models.CreateBannerRequest) error
}

func NewService(logger *slog.Logger, cfg config.Clients, storage *storage.Storage, cash *cash_redis.CashRedis) (*Service, error) {
	const op = "service.NewService"
	triggerHookClient, err := trigger_hook_grpc.NewClient(
		context.Background(),
		logger,
		cfg.TriggerHookService.Address,
		cfg.TriggerHookService.Timeout,
		cfg.TriggerHookService.RetriesCount)
	if err != nil {
		logger.Info("failed connect to trigger-hook-service", fmt.Errorf("%s: %w", op, err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Service{
		Auth:             NewAuthService(storage),
		UserBannerGetter: NewUserBannerManager(storage, cash),
		BannerOperations: NewBannerManager(storage, cash),
		ParamsOperations: NewParamsManager(storage),
		DeleteService:    NewDeleteService(triggerHookClient, storage),
	}, nil
}
