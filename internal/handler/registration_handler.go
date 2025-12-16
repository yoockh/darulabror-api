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

// PUBLIC: POST /registrations
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
		return utils.InternalServerErrorResponse(c, err.Error())
	}

	return c.NoContent(http.StatusCreated)
}

// ADMIN: GET /admin/registrations
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
func (h *RegistrationHandler) AdminGetByID(c echo.Context) error {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, "invalid id")
	}

	item, err := h.svc.GetRegistrationByID(uint(id64))
	if err != nil {
		return utils.NotFoundResponse(c, err.Error())
	}
	return utils.SuccessResponse(c, "registration fetched", item)
}

// ADMIN: DELETE /admin/registrations/:id
func (h *RegistrationHandler) AdminDelete(c echo.Context) error {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, "invalid id")
	}

	if err := h.svc.DeleteRegistration(uint(id64)); err != nil {
		return utils.InternalServerErrorResponse(c, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
