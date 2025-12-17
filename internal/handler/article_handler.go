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
// ListPublished godoc
// @Summary List published articles
// @Description Returns only articles with status "published".
// @Tags Articles (Public)
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Page size" default(10)
// @Success 200 {object} ArticleListResponse
// @Failure 500 {object} ErrorResponse
// @Router /articles [get]
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
// GetPublishedByID godoc
// @Summary Get published article by ID
// @Tags Articles (Public)
// @Produce json
// @Param id path int true "Article ID" minimum(1)
// @Success 200 {object} SuccessResponse[dto.ArticleDTO]
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /articles/{id} [get]
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
// AdminListAll godoc
// @Summary Admin list all articles
// @Description Returns draft + published.
// @Tags Articles (Admin)
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Page size" default(10)
// @Success 200 {object} ArticleListResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/articles [get]
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
// AdminCreate godoc
// @Summary Admin create article
// @Tags Articles (Admin)
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.ArticleDTO true "Article payload (status optional; defaults to draft)"
// @Success 201 {string} string "Created"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/articles [post]
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
// AdminUpdate godoc
// @Summary Admin update article
// @Tags Articles (Admin)
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Article ID" minimum(1)
// @Param request body dto.ArticleDTO true "Article payload (set status=published to publish)"
// @Success 200 {string} string "OK"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/articles/{id} [put]
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
// AdminDelete godoc
// @Summary Admin delete article
// @Tags Articles (Admin)
// @Security BearerAuth
// @Produce json
// @Param id path int true "Article ID" minimum(1)
// @Success 204 {string} string "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/articles/{id} [delete]
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
