package handlers

import (
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

	token, err := h.Service.Auth.Authentication(input)
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

	err := h.Service.Auth.CreateUser(input)
	if err != nil {
		h.Logger.Error("Failed to create user: ", err)
		newErrorResponse(c, http.StatusInternalServerError, internalServerError)
		return
	}

	c.JSON(http.StatusCreated, nil)
}
