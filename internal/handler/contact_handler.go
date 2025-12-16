package handler

import (
	"darulabror/internal/service"
	"darulabror/internal/utils"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type ContactHandler struct {
	svc service.ContactService
}

func NewContactHandler(svc service.ContactService) *ContactHandler {
	return &ContactHandler{svc: svc}
}

// PUBLIC: POST /contacts
func (h *ContactHandler) Create(c echo.Context) error {
	type req struct {
		Email   string `json:"email" validate:"required,email"`
		Subject string `json:"subject" validate:"required,min=3,max=150"`
		Message string `json:"message" validate:"required,min=3,max=2000"`
	}
	var body req

	if err := c.Bind(&body); err != nil {
		return utils.BadRequestResponse(c, "invalid body")
	}
	if err := c.Validate(&body); err != nil {
		return utils.UnprocessableEntityResponse(c, err.Error())
	}

	if err := h.svc.CreateContact(body.Email, body.Subject, body.Message); err != nil {
		logrus.WithError(err).Error("failed create contact")
		return utils.InternalServerErrorResponse(c, "failed to submit contact")
	}

	return c.NoContent(http.StatusCreated)
}

// ADMIN: GET /admin/contacts
func (h *ContactHandler) AdminList(c echo.Context) error {
	page, limit := utils.ParsePagination(c)
	items, total, err := h.svc.GetAllContacts(page, limit)
	if err != nil {
		logrus.WithError(err).Error("failed list contacts")
		return utils.InternalServerErrorResponse(c, "failed to fetch contacts")
	}

	return utils.SuccessResponse(c, "contacts fetched", map[string]interface{}{
		"items": items,
		"meta": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// ADMIN: GET /admin/contacts/:id
func (h *ContactHandler) AdminGetByID(c echo.Context) error {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, "invalid id")
	}

	item, err := h.svc.GetContactByID(uint(id64))
	if err != nil {
		return utils.NotFoundResponse(c, err.Error())
	}
	return utils.SuccessResponse(c, "contact fetched", item)
}

// ADMIN: PUT /admin/contacts/:id
func (h *ContactHandler) AdminUpdate(c echo.Context) error {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, "invalid id")
	}

	type req struct {
		Email   string `json:"email" validate:"required,email"`
		Subject string `json:"subject" validate:"required,min=3,max=150"`
		Message string `json:"message" validate:"required,min=3,max=2000"`
	}
	var body req

	if err := c.Bind(&body); err != nil {
		return utils.BadRequestResponse(c, "invalid body")
	}
	if err := c.Validate(&body); err != nil {
		return utils.UnprocessableEntityResponse(c, err.Error())
	}

	if err := h.svc.UpdateContact(uint(id64), body.Email, body.Subject, body.Message); err != nil {
		return utils.InternalServerErrorResponse(c, err.Error())
	}
	return c.NoContent(http.StatusOK)
}

// ADMIN: DELETE /admin/contacts/:id
func (h *ContactHandler) AdminDelete(c echo.Context) error {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, "invalid id")
	}

	if err := h.svc.DeleteContact(uint(id64)); err != nil {
		return utils.InternalServerErrorResponse(c, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
