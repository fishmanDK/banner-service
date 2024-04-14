package handlers

import (
	"fmt"
	"github.com/fishmanDK/avito_test_task/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handlers) CreateTag(c *gin.Context) {
	const op = "handlers.CreateTag"

	var input models.Tag
	if err := c.BindJSON(&input); err != nil {
		h.Logger.Error("Failed to bind JSON: ", fmt.Errorf("%s: %v", op, err))
		newErrorResponse(c, http.StatusBadRequest, incorrectData)
		return
	}

	err := h.Service.ParamsOperations.CreateTag(input)
	if err != nil {
		h.Logger.Error("Error creating banner: ", fmt.Errorf("%s: %v", op, err))
		newErrorResponse(c, http.StatusInternalServerError, internalServerError)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *Handlers) CreateFeature(c *gin.Context) {
	const op = "handlers.CreateBanner"

	var input models.Feature
	if err := c.BindJSON(&input); err != nil {
		h.Logger.Error("Failed to bind JSON: ", fmt.Errorf("%s: %v", op, err))
		newErrorResponse(c, http.StatusBadRequest, incorrectData)
		return
	}

	err := h.Service.ParamsOperations.CreateFeature(input)
	if err != nil {
		h.Logger.Error("Error creating banner: ", fmt.Errorf("%s: %v", op, err))
		newErrorResponse(c, http.StatusInternalServerError, internalServerError)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
