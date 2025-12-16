package handler

import (
	"darulabror/internal/dto"
	"darulabror/internal/models"
	"darulabror/internal/service"
	"darulabror/internal/utils"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type AdminHandler struct {
	svc service.AdminService
}

func NewAdminHandler(svc service.AdminService) *AdminHandler {
	return &AdminHandler{svc: svc}
}

// SUPERADMIN: POST /admin/admins
func (h *AdminHandler) Create(c echo.Context) error {
	var body dto.AdminDTO
	if err := c.Bind(&body); err != nil {
		return utils.BadRequestResponse(c, "invalid body")
	}
	if err := c.Validate(&body); err != nil {
		return utils.UnprocessableEntityResponse(c, err.Error())
	}
	if body.Password == "" {
		return utils.UnprocessableEntityResponse(c, "password is required")
	}

	role, _ := utils.GetRole(c)
	if err := h.svc.CreateAdmin(role, body); err != nil {
		if err.Error() == "forbidden" {
			return utils.ForbiddenResponse(c, "forbidden")
		}
		return utils.InternalServerErrorResponse(c, err.Error())
	}

	return c.NoContent(http.StatusCreated)
}

// SUPERADMIN: GET /admin/admins
func (h *AdminHandler) List(c echo.Context) error {
	page, limit := utils.ParsePagination(c)
	items, total, err := h.svc.GetAllAdmins(page, limit)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "failed to fetch admins")
	}

	return utils.SuccessResponse(c, "admins fetched", map[string]interface{}{
		"items": items,
		"meta": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// ADMIN/SUPERADMIN: GET /admin/profile (contoh simpel: get by id from token)
func (h *AdminHandler) Profile(c echo.Context) error {
	adminID, ok := utils.GetAdminID(c)
	if !ok {
		return utils.UnauthorizedResponse(c, "unauthorized")
	}
	item, err := h.svc.GetAdminByID(adminID)
	if err != nil {
		return utils.NotFoundResponse(c, "admin not found")
	}
	return utils.SuccessResponse(c, "profile fetched", item)
}

// SUPERADMIN: PUT /admin/admins/:id
func (h *AdminHandler) Update(c echo.Context) error {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, "invalid id")
	}

	var body dto.AdminDTO
	if err := c.Bind(&body); err != nil {
		return utils.BadRequestResponse(c, "invalid body")
	}
	if err := c.Validate(&body); err != nil {
		return utils.UnprocessableEntityResponse(c, err.Error())
	}
	body.ID = uint(id64)

	role, _ := utils.GetRole(c)
	if err := h.svc.UpdateAdmin(role, body); err != nil {
		if err.Error() == "forbidden" {
			return utils.ForbiddenResponse(c, "forbidden")
		}
		return utils.InternalServerErrorResponse(c, err.Error())
	}
	return c.NoContent(http.StatusOK)
}

// SUPERADMIN: DELETE /admin/admins/:id
func (h *AdminHandler) Delete(c echo.Context) error {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, "invalid id")
	}

	role, _ := utils.GetRole(c)
	if err := h.svc.DeleteAdmin(role, uint(id64)); err != nil {
		if err.Error() == "forbidden" {
			return utils.ForbiddenResponse(c, "forbidden")
		}
		return utils.InternalServerErrorResponse(c, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

// (optional) helper supaya compile kalau dipakai di routes
func allowAdminOrSuperadmin(_ models.Role) bool { return true }
