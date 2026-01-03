package handler

import (
	"context"
	"darulabror/internal/dto"
	"darulabror/internal/repository"
	"darulabror/internal/service"
	"darulabror/internal/utils"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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
// @Summary Admin create article (multipart)
// @Tags Articles (Admin)
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param title formData string true "Title"
// @Param author formData string true "Author"
// @Param status formData string false "draft|published" Enums(draft,published)
// @Param content formData string true "JSON string (flexible)"
// @Param photo_header formData string false "Optional header URL (ignored if photo_header_file is provided)"
// @Param photo_header_file formData file false "Optional header image file (uploaded and set to photo_header)"
// @Param content_files formData file false "Inline media files. Use field name: content_files[<upload_key>] (repeatable). Example: content_files[img1], content_files[vid1]"
// @Success 201 {string} string "Created"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/articles [post]
func (h *ArticleHandler) AdminCreate(c echo.Context) error {
	title := c.FormValue("title")
	author := c.FormValue("author")
	status := c.FormValue("status")
	contentStr := c.FormValue("content")

	if title == "" || author == "" || contentStr == "" {
		return utils.BadRequestResponse(c, "missing required fields: title, author, content")
	}

	// 1) parse content JSON
	if !json.Valid([]byte(contentStr)) {
		return utils.BadRequestResponse(c, "content must be valid JSON string")
	}
	var contentAny any
	if err := json.Unmarshal([]byte(contentStr), &contentAny); err != nil {
		return utils.BadRequestResponse(c, "content must be valid JSON")
	}

	// 2) upload inline content files (content_files[<key>])
	urlByKey, err := h.parseAndUploadContentFiles(c)
	if err != nil {
		if errors.Is(err, repository.ErrStorageNotConfigured) {
			return utils.BadRequestResponse(c, "storage not configured: set PUBLIC_BUCKET to enable uploads")
		}
		logrus.WithError(err).Error("failed upload content files")
		return utils.InternalServerErrorResponse(c, "failed to upload content files")
	}

	// 3) inject URLs into content JSON
	if len(urlByKey) > 0 {
		contentAny = injectUploadedURLs(contentAny, urlByKey)
	}
	contentBytes, err := json.Marshal(contentAny)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "failed to serialize content")
	}

	body := dto.ArticleDTO{
		Title:       title,
		Author:      author,
		Status:      status,
		PhotoHeader: c.FormValue("photo_header"),
		Content:     contentBytes,
	}

	// 4) optional header upload (photo_header_file) overrides photo_header
	if fh, err := c.FormFile("photo_header_file"); err == nil && fh != nil {
		src, err := fh.Open()
		if err != nil {
			return utils.BadRequestResponse(c, "failed to open photo_header_file")
		}
		defer src.Close()

		safeName := filepath.Base(fh.Filename)
		objectName := "articles/header_" + strconv.FormatInt(time.Now().UnixNano(), 10) + "_" + safeName

		urlOrObject, err := h.svc.UploadArticleMedia(c.Request().Context(), src, objectName)
		if err != nil {
			if errors.Is(err, repository.ErrStorageNotConfigured) {
				return utils.BadRequestResponse(c, "storage not configured: set PUBLIC_BUCKET to enable uploads")
			}
			logrus.WithError(err).Error("failed upload photo_header_file")
			return utils.InternalServerErrorResponse(c, "failed to upload header")
		}
		body.PhotoHeader = urlOrObject
	}

	// after parsing fields + (optional) uploading photo_header_file
	if strings.TrimSpace(body.PhotoHeader) == "" {
		return utils.BadRequestResponse(c, "photo_header is required (provide photo_header URL or upload photo_header_file)")
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
// @Summary Admin update article (multipart)
// @Tags Articles (Admin)
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Article ID" minimum(1)
// @Param title formData string true "Title"
// @Param author formData string true "Author"
// @Param status formData string false "draft|published" Enums(draft,published)
// @Param content formData string true "JSON string (flexible)"
// @Param photo_header formData string false "Optional header URL (ignored if photo_header_file is provided)"
// @Param photo_header_file formData file false "Optional header image file (uploaded and set to photo_header)"
// @Param content_files formData file false "Inline media files. Use field name: content_files[<upload_key>] (repeatable). Example: content_files[img1], content_files[vid1]"
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

	title := c.FormValue("title")
	author := c.FormValue("author")
	status := c.FormValue("status")
	contentStr := c.FormValue("content")

	if title == "" || author == "" || contentStr == "" {
		return utils.BadRequestResponse(c, "missing required fields: title, author, content")
	}

	// Same content flow as create:
	if !json.Valid([]byte(contentStr)) {
		return utils.BadRequestResponse(c, "content must be valid JSON string")
	}
	var contentAny any
	if err := json.Unmarshal([]byte(contentStr), &contentAny); err != nil {
		return utils.BadRequestResponse(c, "content must be valid JSON")
	}

	urlByKey, err := h.parseAndUploadContentFiles(c)
	if err != nil {
		if errors.Is(err, repository.ErrStorageNotConfigured) {
			return utils.BadRequestResponse(c, "storage not configured: set PUBLIC_BUCKET to enable uploads")
		}
		logrus.WithError(err).Error("failed upload content files")
		return utils.InternalServerErrorResponse(c, "failed to upload content files")
	}
	if len(urlByKey) > 0 {
		contentAny = injectUploadedURLs(contentAny, urlByKey)
	}
	contentBytes, err := json.Marshal(contentAny)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "failed to serialize content")
	}

	body := dto.ArticleDTO{
		Title:       title,
		Author:      author,
		Status:      status,
		PhotoHeader: c.FormValue("photo_header"),
		Content:     contentBytes,
	}

	// Optional header upload: if provided, overrides photo_header string
	if fh, err := c.FormFile("photo_header_file"); err == nil && fh != nil {
		src, err := fh.Open()
		if err != nil {
			return utils.BadRequestResponse(c, "failed to open photo_header_file")
		}
		defer src.Close()

		safeName := filepath.Base(fh.Filename)
		objectName := "articles/header_" + strconv.FormatInt(time.Now().UnixNano(), 10) + "_" + safeName

		urlOrObject, err := h.svc.UploadArticleMedia(c.Request().Context(), src, objectName)
		if err != nil {
			if errors.Is(err, repository.ErrStorageNotConfigured) {
				return utils.BadRequestResponse(c, "storage not configured: set PUBLIC_BUCKET to enable header upload")
			}
			logrus.WithError(err).Error("failed upload photo_header_file")
			return utils.InternalServerErrorResponse(c, "failed to upload header")
		}
		body.PhotoHeader = urlOrObject
	}

	// after parsing fields + (optional) uploading photo_header_file
	if strings.TrimSpace(body.PhotoHeader) == "" {
		return utils.BadRequestResponse(c, "photo_header is required (provide photo_header URL or upload photo_header_file)")
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
		logrus.WithError(err).WithField("id", id64).Error("failed delete article")
		return utils.InternalServerErrorResponse(c, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

// extractUploadKey supports field naming:
// - content_files[img1]
// - content_file_img1 (fallback)
func extractUploadKey(field string) (string, bool) {
	if strings.HasPrefix(field, "content_files[") && strings.HasSuffix(field, "]") {
		key := strings.TrimSuffix(strings.TrimPrefix(field, "content_files["), "]")
		key = strings.TrimSpace(key)
		return key, key != ""
	}
	if strings.HasPrefix(field, "content_file_") {
		key := strings.TrimSpace(strings.TrimPrefix(field, "content_file_"))
		return key, key != ""
	}
	return "", false
}

func uploadOne(ctx context.Context, svc interface {
	UploadArticleMedia(ctx context.Context, file io.Reader, objectName string) (string, error)
}, fh *multipart.FileHeader, objectPrefix, key string) (string, error) {
	src, err := fh.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	safeName := filepath.Base(fh.Filename)
	objectName := objectPrefix + "/" + key + "_" + strconv.FormatInt(time.Now().UnixNano(), 10) + "_" + safeName
	return svc.UploadArticleMedia(ctx, src, objectName)
}

// replace any map that contains {"upload_key": "<key>"} with {"url": "<url>"} (and removes upload_key)
// works recursively for objects/arrays
func injectUploadedURLs(node any, urlByKey map[string]string) any {
    switch v := node.(type) {
    case map[string]any:
        // Case 1: legacy placeholder { upload_key: "<key>" }
        if raw, ok := v["upload_key"]; ok {
            if key, ok := raw.(string); ok {
                if url, exists := urlByKey[key]; exists {
                    v["url"] = url
                    delete(v, "upload_key")
                }
            }
        }

        // Case 2: EditorJS image block: {type:"image", data:{file:{fileKey, url}}}
        if t, ok := v["type"].(string); ok && strings.EqualFold(t, "image") {
            if data, ok := v["data"].(map[string]any); ok {
                if file, ok := data["file"].(map[string]any); ok {
                    // support several key names just in case
                    var key string
                    if k, ok := file["fileKey"].(string); ok && strings.TrimSpace(k) != "" {
                        key = strings.TrimSpace(k)
                    } else if k, ok := file["file_key"].(string); ok && strings.TrimSpace(k) != "" {
                        key = strings.TrimSpace(k)
                    }

                    if key != "" {
                        if url, exists := urlByKey[key]; exists && strings.TrimSpace(url) != "" {
                            file["url"] = url // overwrite blob: with public https URL
                            // optional: you may keep fileKey for future edits, or delete it:
                            // delete(file, "fileKey")
                        }
                    }
                }
            }
        }

        // Recurse through children
        for k, child := range v {
            v[k] = injectUploadedURLs(child, urlByKey)
        }
        return v

    case []any:
        for i := range v {
            v[i] = injectUploadedURLs(v[i], urlByKey)
        }
        return v

    default:
        return node
    }
}

func (h *ArticleHandler) parseAndUploadContentFiles(c echo.Context) (map[string]string, error) {
	form, err := c.MultipartForm()
	if err != nil {
		// no multipart form or not parsed: treat as no files
		return map[string]string{}, nil
	}

	urlByKey := map[string]string{}
	for field, fhs := range form.File {
		key, ok := extractUploadKey(field)
		if !ok {
			continue
		}
		if len(fhs) == 0 {
			continue
		}

		// Only first file per key (keep API predictable)
		url, err := uploadOne(c.Request().Context(), h.svc, fhs[0], "articles/content", key)
		if err != nil {
			if errors.Is(err, repository.ErrStorageNotConfigured) {
				return nil, repository.ErrStorageNotConfigured
			}
			return nil, err
		}
		urlByKey[key] = url
	}

	return urlByKey, nil
}
