package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/fishmanDK/avito_test_task/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

const (
	bannerNotSelected = "Баннер не выбран"
	admin             = "admin"
	ordinary          = "ordinary"
)

func (h *Handlers) GetBanners(c *gin.Context) {
	const op = "handlers.GetBanners"

	var params models.GetAllBannersParams

	featureIDStr := c.Query("feature_id")
	tagIDStr := c.Query("tag_id")
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")

	if featureIDStr != "" {
		featureID, err := strconv.Atoi(featureIDStr)
		if err != nil {
			h.Logger.Error("Failed get feature_id: ", fmt.Errorf("%s: %v", op, err))
			newErrorResponse(c, http.StatusBadRequest, incorrectData)
			return
		}
		featureIDPtr := &featureID
		params.FeatureID = featureIDPtr
	}

	if tagIDStr != "" {
		tagID, err := strconv.Atoi(tagIDStr)
		if err != nil {
			h.Logger.Error("Failed get tag_id: ", fmt.Errorf("%s: %v", op, err))
			newErrorResponse(c, http.StatusBadRequest, incorrectData)
			return
		}
		tagIDPtr := &tagID
		params.TagID = tagIDPtr
	}

	if limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			h.Logger.Error("Failed get limit: ", fmt.Errorf("%s: %v", op, err))
			newErrorResponse(c, http.StatusBadRequest, incorrectData)
			return
		}
		limitPtr := &limit
		params.Limit = limitPtr
	}

	if offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			h.Logger.Error("Failed get offset: ", fmt.Errorf("%s: %v", op, err))
			newErrorResponse(c, http.StatusBadRequest, incorrectData)
			return
		}
		offsetPtr := &offset
		params.Offset = offsetPtr
	}

	ctx, cancel := context.WithTimeout(context.Background(), contextTimeResponse)
	defer cancel()

	banners, err := h.Service.GetBanners(ctx, params)
	if err != nil {
		h.Logger.Error("Error getting banners: ", fmt.Errorf("%s: %v", op, err))
		newErrorResponse(c, http.StatusBadRequest, incorrectData)
		return
	}

	c.JSON(http.StatusOK, banners)
}

func (h *Handlers) CreateBanner(c *gin.Context) {
	const op = "handlers.CreateBanner"

	var input models.CreateBannerRequest
	if err := c.BindJSON(&input); err != nil {
		h.Logger.Error("Failed to bind JSON: ", fmt.Errorf("%s: %v", op, err))
		newErrorResponse(c, http.StatusBadRequest, incorrectData)
		return
	}

	_, err := strconv.Atoi(input.FeatureID)
	if err != nil {
		h.Logger.Error("Incorrect feature_id: ", fmt.Errorf("%s: %v", op, err))
		newErrorResponse(c, http.StatusBadRequest, incorrectData)
		return
	}

	input.IsActive = strings.ToUpper(input.IsActive)
	if !strings.EqualFold(input.IsActive, "TRUE") && !strings.EqualFold(input.IsActive, "FALSE") {
		h.Logger.Error("Incorrect is_active: ", fmt.Errorf("%s: %v", op, err))
		newErrorResponse(c, http.StatusBadRequest, incorrectData)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeResponse)
	defer cancel()

	err = h.Service.BannerOperations.CreateBanner(ctx, input)
	if err != nil {
		h.Logger.Error("Error creating banner: ", fmt.Errorf("%s: %v", op, err))
		newErrorResponse(c, http.StatusInternalServerError, internalServerError)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *Handlers) PatchBanner(c *gin.Context) {
	const op = "handlers.PatchBanner"

	bannerIDStr := c.Param("id")
	if bannerIDStr == "" {
		h.Logger.Error("Incorrect data: ", fmt.Errorf("%s: %v", op, errors.New("ID parameter is missing")))
		newErrorResponse(c, http.StatusBadRequest, bannerNotSelected)
		return
	}

	bannerID, err := strconv.ParseInt(bannerIDStr, 10, 64)
	if err != nil {
		h.Logger.Error("Failed convert bannerID to int64: ", fmt.Errorf("%s: %v", op, err))
		newErrorResponse(c, http.StatusBadRequest, incorrectData)
		return
	}

	var input models.ChangeBannerRequest
	if err := c.BindJSON(&input); err != nil {
		h.Logger.Error("Failed to bind JSON: ", fmt.Errorf("%s: %v", op, err))
		newErrorResponse(c, http.StatusBadRequest, incorrectData)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), contextTimeResponse)
	defer cancel()

	err = h.Service.BannerOperations.ChangeBanner(ctx, bannerID, input)
	if err != nil {
		h.Logger.Error("Error changing banner: ", err)
		newErrorResponse(c, http.StatusInternalServerError, internalServerError)
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (h *Handlers) DeleteBanner(c *gin.Context) {
	const op = "handlers.DeleteBanner"

	bannerIDStr := c.Param("id")
	if bannerIDStr == "" {
		h.Logger.Error("Incorrect data: ", fmt.Errorf("%s: %v", op, errors.New("ID parameter is missing")))
		newErrorResponse(c, http.StatusBadRequest, bannerNotSelected)
		return
	}

	bannerID, err := strconv.ParseInt(bannerIDStr, 10, 64)
	if err != nil {
		h.Logger.Error("Failed convert bannerID to int64: ", fmt.Errorf("%s: %v", op, err))
		newErrorResponse(c, http.StatusBadRequest, incorrectData)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), contextTimeResponse)
	defer cancel()

	err = h.Service.DeleteService.ScheduleDeleteBanner(ctx, bannerID, 0, 0)
	if err != nil {
		h.Logger.Error("Error deleting banner: ", fmt.Errorf("%s: %v", op, err))
		newErrorResponse(c, http.StatusBadRequest, incorrectData)
		return
	}

	c.JSON(http.StatusAccepted, nil)
}

func (h *Handlers) DeleteBannerByParams(c *gin.Context) {
	const op = "handlers.DeleteBannerByParams"

	tagIDStr, ok := c.GetQuery("tag_id")
	if !ok {
		newErrorResponse(c, http.StatusBadRequest, incorrectData)
		return
	}

	tagID, err := strconv.ParseInt(tagIDStr, 10, 64)
	if err != nil {
		h.Logger.Error("Failed convert bannerID to int64: ", fmt.Errorf("%s: %v", op, err))
		newErrorResponse(c, http.StatusBadRequest, incorrectData)
		return
	}

	featureIDStr, ok := c.GetQuery("feature_id")
	if !ok {
		newErrorResponse(c, http.StatusBadRequest, incorrectData)
		return
	}

	featureID, err := strconv.ParseInt(featureIDStr, 10, 64)
	if err != nil {
		h.Logger.Error("Failed convert bannerID to int64: ", fmt.Errorf("%s: %v", op, err))
		newErrorResponse(c, http.StatusBadRequest, incorrectData)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), contextTimeResponse)
	defer cancel()

	err = h.Service.DeleteService.ScheduleDeleteBanner(ctx, 0, tagID, featureID)
	if err != nil {
		h.Logger.Error("Error deleting banner: ", fmt.Errorf("%s: %v", op, err))
		newErrorResponse(c, http.StatusBadRequest, incorrectData)
		return
	}

	c.JSON(http.StatusAccepted, nil)
}
