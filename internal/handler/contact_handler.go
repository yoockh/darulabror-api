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

// Create godoc
// @Summary Create contact message
// @Tags Contacts (Public)
// @Accept json
// @Produce json
// @Param request body ContactCreateRequest true "Contact payload"
// @Success 201 {string} string "Created"
// @Failure 400 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /contacts [post]
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
// AdminList godoc
// @Summary Admin list contacts
// @Tags Contacts (Admin)
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Page size" default(10)
// @Success 200 {object} ContactListResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/contacts [get]
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
// AdminGetByID godoc
// @Summary Admin get contact by ID
// @Tags Contacts (Admin)
// @Security BearerAuth
// @Produce json
// @Param id path int true "Contact ID" minimum(1)
// @Success 200 {object} SuccessResponse[models.Contact]
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /admin/contacts/{id} [get]
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
// AdminUpdate godoc
// @Summary Admin update contact
// @Tags Contacts (Admin)
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Contact ID" minimum(1)
// @Param request body ContactUpdateRequest true "Contact payload"
// @Success 200 {string} string "OK"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/contacts/{id} [put]
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
// AdminDelete godoc
// @Summary Admin delete contact
// @Tags Contacts (Admin)
// @Security BearerAuth
// @Produce json
// @Param id path int true "Contact ID" minimum(1)
// @Success 204 {string} string "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/contacts/{id} [delete]
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
