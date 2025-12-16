package utils

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

const (
	DefaultPage  = 1
	DefaultLimit = 10
	MaxLimit     = 100
)

// NormalizePageLimit returns (page, limit, offset)
func NormalizePageLimit(page, limit int) (int, int, int) {
	if page <= 0 {
		page = DefaultPage
	}
	if limit <= 0 {
		limit = DefaultLimit
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}
	offset := (page - 1) * limit
	return page, limit, offset
}

func ParsePagination(c echo.Context) (page, limit int) {
	page, _ = strconv.Atoi(c.QueryParam("page"))
	limit, _ = strconv.Atoi(c.QueryParam("limit"))
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	return
}
