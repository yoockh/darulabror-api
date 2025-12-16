package handler

import (
	"darulabror/internal/dto"
	"darulabror/internal/service"
	"darulabror/internal/utils"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type ArticleHandler struct {
	svc service.ArticleService
}

func NewArticleHandler(svc service.ArticleService) *ArticleHandler {
	return &ArticleHandler{svc: svc}
}

// PUBLIC: GET /articles
func (h *ArticleHandler) ListPublished(c echo.Context) error {
	page, limit := utils.ParsePagination(c)
	items, total, err := h.svc.GetPublishedArticles(page, limit)
	if err != nil {
		logrus.WithError(err).Error("failed list published articles")
		return utils.InternalServerErrorResponse(c, "failed to fetch articles")
	}

	return utils.SuccessResponse(c, "articles fetched", map[string]interface{}{
		"items": items,
		"meta": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// PUBLIC: GET /articles/:id
func (h *ArticleHandler) GetPublishedByID(c echo.Context) error {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, "invalid id")
	}

	item, err := h.svc.GetPublishedArticleByID(uint(id64))
	if err != nil {
		return utils.NotFoundResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, "article fetched", item)
}

// ADMIN: GET /admin/articles
func (h *ArticleHandler) AdminListAll(c echo.Context) error {
	page, limit := utils.ParsePagination(c)
	items, total, err := h.svc.GetAllArticles(page, limit)
	if err != nil {
		logrus.WithError(err).Error("failed admin list all articles")
		return utils.InternalServerErrorResponse(c, "failed to fetch articles")
	}

	return utils.SuccessResponse(c, "articles fetched", map[string]interface{}{
		"items": items,
		"meta": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// ADMIN: POST /admin/articles
func (h *ArticleHandler) AdminCreate(c echo.Context) error {
	var body dto.ArticleDTO
	if err := c.Bind(&body); err != nil {
		return utils.BadRequestResponse(c, "invalid body")
	}
	if err := c.Validate(&body); err != nil {
		return utils.UnprocessableEntityResponse(c, err.Error())
	}

	if err := h.svc.CreateArticle(body); err != nil {
		return utils.InternalServerErrorResponse(c, err.Error())
	}
	return c.NoContent(http.StatusCreated)
}

// ADMIN: PUT /admin/articles/:id
func (h *ArticleHandler) AdminUpdate(c echo.Context) error {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, "invalid id")
	}

	var body dto.ArticleDTO
	if err := c.Bind(&body); err != nil {
		return utils.BadRequestResponse(c, "invalid body")
	}
	if err := c.Validate(&body); err != nil {
		return utils.UnprocessableEntityResponse(c, err.Error())
	}

	if err := h.svc.UpdateArticle(uint(id64), body); err != nil {
		return utils.InternalServerErrorResponse(c, err.Error())
	}
	return c.NoContent(http.StatusOK)
}

// ADMIN: DELETE /admin/articles/:id
func (h *ArticleHandler) AdminDelete(c echo.Context) error {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, "invalid id")
	}

	if err := h.svc.DeleteArticle(uint(id64)); err != nil {
		return utils.InternalServerErrorResponse(c, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
