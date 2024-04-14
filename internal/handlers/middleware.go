package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	errorAuthorize  = "ошибка авторизации"
	errorParseToken = "ошибка парсинга токена"
)

func (h *Handlers) authMiddleware(c *gin.Context) {
	const op = "handlers.authMiddleware"

	header := c.GetHeader("Authorization")
	if header == "" {
		h.Logger.Error("Error get header (Authorization): ", op)
		newErrorResponse(c, http.StatusUnauthorized, userUnauthorized)
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		h.Logger.Error("Error format header (Authorization): ", op)
		newErrorResponse(c, http.StatusUnauthorized, errorAuthorize)
		return
	}

	role, err := h.Service.Auth.ParseToken(headerParts[1])
	if err != nil {
		h.Logger.Error("Error parse token: ", op)
		newErrorResponse(c, http.StatusUnauthorized, errorParseToken)
		return
	}

	c.Set("role", role)
	c.Next()
}

func (h *Handlers) isAdminMiddleware(c *gin.Context) {
	const op = "handlers.isAdminMiddleware"

	role, ok := c.Get("role")
	if !ok {
		h.Logger.Error("Error get role: ", op)
		newErrorResponse(c, http.StatusUnauthorized, userUnauthorized)
		return
	}

	if role != admin {
		h.Logger.Error("User is not admin: ", op)
		newErrorResponse(c, http.StatusForbidden, userAccessDenied)
		return
	}
	c.Next()
}
