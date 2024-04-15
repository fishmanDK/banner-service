package handlers

import (
	"github.com/fishmanDK/avito_test_task/internal/service"
	"github.com/gin-gonic/gin"
	"log/slog"
	"time"
)

const (
	contextTimeResponse = 50 * time.Millisecond
)

type Handlers struct {
	Service *service.Service
	Logger  *slog.Logger
}

func MustHandlers(service *service.Service, logger *slog.Logger) *Handlers {
	return &Handlers{
		Service: service,
		Logger:  logger,
	}
}

func (h *Handlers) InitRouts() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.GET("/sign-in", h.signIn)
		auth.POST("/sign-up", h.signUp)
	}

	router.Use(h.authMiddleware)
	router.GET("/user_banner", h.GetUserBanner)

	router.POST("/tag", h.CreateTag)
	router.POST("/feature", h.CreateFeature)

	banner := router.Group("/banner")
	banner.Use(h.isAdminMiddleware)
	{
		banner.GET("", h.GetBanners)
		banner.POST("", h.CreateBanner)
		banner.DELETE("", h.DeleteBannerByParams)

		banner.PATCH("/:id", h.PatchBanner)
		banner.DELETE("/:id", h.DeleteBanner)
	}

	return router
}
