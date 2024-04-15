package handlers

import (
	"context"
	"github.com/fishmanDK/avito_test_task/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	incorrectData       = "Некорректные данные"
	userUnauthorized    = "Пользователь не авторизован"
	userAccessDenied    = "Пользователь не имеет доступа"
	bannerNotFound      = "Баннер не найден"
	internalServerError = "Внутренняя ошибка сервера"
)

func (h *Handlers) GetUserBanner(c *gin.Context) {
	tag_id, ok := c.GetQuery("tag_id")
	if !ok {
		newErrorResponse(c, http.StatusBadRequest, incorrectData)
		return
	}

	feature_id, ok := c.GetQuery("feature_id")
	if !ok {
		newErrorResponse(c, http.StatusBadRequest, incorrectData)
		return
	}

	req := models.UserBanner{
		TagID:     tag_id,
		FeatureID: feature_id,
	}

	use_last_revision, ok := c.GetQuery("use_last_revision")
	if use_last_revision != "true" {
		req.UseLastRevision = false
	} else {
		req.UseLastRevision = true
	}

	ctx, cancel := context.WithTimeout(context.Background(), contextTimeResponse)
	defer cancel()

	banner, err := h.Service.GetUserBanner(ctx, req)
	if err != nil {
		h.Logger.Error("error storage", err)
		newErrorResponse(c, http.StatusBadRequest, incorrectData)
		return
	}

	c.JSON(http.StatusOK, banner)
}
