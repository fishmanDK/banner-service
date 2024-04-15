package handlers

import (
	"context"
	"github.com/fishmanDK/avito_test_task/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handlers) signIn(c *gin.Context) {
	var input models.User
	if err := c.BindJSON(&input); err != nil {
		h.Logger.Error("Failed to bind JSON: ", err)
		newErrorResponse(c, http.StatusBadRequest, incorrectData)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), contextTimeResponse)
	defer cancel()

	token, err := h.Service.Auth.Authentication(ctx, input)
	if err != nil {
		h.Logger.Error("Authentication failed: ", err)
		newErrorResponse(c, http.StatusInternalServerError, internalServerError)
		return
	}

	c.JSON(http.StatusOK, token)
}

func (h *Handlers) signUp(c *gin.Context) {
	var input models.NewUser
	if err := c.BindJSON(&input); err != nil {
		h.Logger.Error("Failed to bind JSON: ", err)
		newErrorResponse(c, http.StatusBadRequest, incorrectData)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeResponse)
	defer cancel()

	err := h.Service.Auth.CreateUser(ctx, input)
	if err != nil {
		h.Logger.Error("Failed to create user: ", err)
		newErrorResponse(c, http.StatusInternalServerError, internalServerError)
		return
	}

	c.JSON(http.StatusCreated, nil)
}
