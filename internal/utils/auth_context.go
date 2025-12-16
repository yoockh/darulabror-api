package utils

import (
	"darulabror/internal/models"

	"github.com/labstack/echo/v4"
)

const (
	CtxAdminIDKey = "admin_id"
	CtxRoleKey    = "role"
)

func GetAdminID(c echo.Context) (uint, bool) {
	v := c.Get(CtxAdminIDKey)
	if v == nil {
		return 0, false
	}
	id, ok := v.(uint)
	return id, ok
}

func GetRole(c echo.Context) (models.Role, bool) {
	v := c.Get(CtxRoleKey)
	if v == nil {
		return "", false
	}
	role, ok := v.(models.Role)
	return role, ok
}
