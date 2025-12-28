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

type RegistrationHandler struct {
	svc service.RegistrationService
}

func NewRegistrationHandler(svc service.RegistrationService) *RegistrationHandler {
	return &RegistrationHandler{svc: svc}
}

// Create godoc
// @Summary Create registration
// @Tags Registrations (Public)
// @Accept json
// @Produce json
// @Param request body dto.RegistrationDTO true "Registration payload"
// @Success 201 {string} string "Created"
// @Failure 400 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /registrations [post]
func (h *RegistrationHandler) Create(c echo.Context) error {
	var body dto.RegistrationDTO
	if err := c.Bind(&body); err != nil {
		return utils.BadRequestResponse(c, "invalid body")
	}
	if err := c.Validate(&body); err != nil {
		return utils.UnprocessableEntityResponse(c, err.Error())
	}

	if err := h.svc.CreateRegistration(body); err != nil {
		logrus.WithError(err).Error("failed create registration")
		return utils.InternalServerErrorResponse(c, "failed to process registration")
	}

	return c.NoContent(http.StatusCreated)
}

// ADMIN: GET /admin/registrations
// AdminList godoc
// @Summary Admin list registrations
// @Tags Registrations (Admin)
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Page size" default(10)
// @Success 200 {object} RegistrationListResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/registrations [get]
func (h *RegistrationHandler) AdminList(c echo.Context) error {
	page, limit := utils.ParsePagination(c)
	items, total, err := h.svc.GetAllRegistrations(page, limit)
	if err != nil {
		logrus.WithError(err).Error("failed list registrations")
		return utils.InternalServerErrorResponse(c, "failed to fetch registrations")
	}

	return utils.SuccessResponse(c, "registrations fetched", map[string]interface{}{
		"items": items,
		"meta": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// ADMIN: GET /admin/registrations/:id
// AdminGetByID godoc
// @Summary Admin get registration by ID
// @Tags Registrations (Admin)
// @Security BearerAuth
// @Produce json
// @Param id path int true "Registration ID" minimum(1)
// @Success 200 {object} SuccessResponse[dto.RegistrationDTO]
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /admin/registrations/{id} [get]
func (h *RegistrationHandler) AdminGetByID(c echo.Context) error {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, "invalid id")
	}

	item, err := h.svc.GetRegistrationByID(uint(id64))
	if err != nil {
		logrus.WithError(err).Error("failed get registration by id")
		return utils.NotFoundResponse(c, "registration not found")
	}
	return utils.SuccessResponse(c, "registration fetched", item)
}

// ADMIN: DELETE /admin/registrations/:id
// AdminDelete godoc
// @Summary Admin delete registration
// @Tags Registrations (Admin)
// @Security BearerAuth
// @Produce json
// @Param id path int true "Registration ID" minimum(1)
// @Success 204 {string} string "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/registrations/{id} [delete]
func (h *RegistrationHandler) AdminDelete(c echo.Context) error {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, "invalid id")
	}

	if err := h.svc.DeleteRegistration(uint(id64)); err != nil {
		logrus.WithError(err).Error("failed delete registration")
		return utils.InternalServerErrorResponse(c, "failed to delete registration")
	}
	return c.NoContent(http.StatusNoContent)
}
