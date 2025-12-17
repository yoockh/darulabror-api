package handler

import (
	"darulabror/internal/dto"
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

// Create godoc
// @Summary Superadmin create admin
// @Tags Admins (Superadmin)
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.AdminDTO true "Admin payload (password required)"
// @Success 201 {string} string "Created"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/admins [post]
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

// List godoc
// @Summary Superadmin list admins
// @Tags Admins (Superadmin)
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Page size" default(10)
// @Success 200 {object} AdminListResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/admins [get]
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

// Profile godoc
// @Summary Get admin profile (from JWT)
// @Tags Admins (Admin)
// @Security BearerAuth
// @Produce json
// @Success 200 {object} SuccessResponse[dto.AdminDTO]
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /admin/profile [get]
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

// Update godoc
// @Summary Superadmin update admin
// @Tags Admins (Superadmin)
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Admin ID" minimum(1)
// @Param request body dto.AdminDTO true "Admin payload (password optional)"
// @Success 200 {string} string "OK"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/admins/{id} [put]
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

// Delete godoc
// @Summary Superadmin delete admin
// @Tags Admins (Superadmin)
// @Security BearerAuth
// @Produce json
// @Param id path int true "Admin ID" minimum(1)
// @Success 204 {string} string "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/admins/{id} [delete]
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

// Login godoc
// @Summary Admin login
// @Description Returns JWT token for accessing /admin endpoints.
// @Tags Auth (Admin)
// @Accept json
// @Produce json
// @Param request body AdminLoginRequest true "Login payload"
// @Success 200 {object} AdminLoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/login [post]
func (h *AdminHandler) Login(c echo.Context) error {
	var body AdminLoginRequest

	if err := c.Bind(&body); err != nil {
		return utils.BadRequestResponse(c, "invalid body")
	}
	if err := c.Validate(&body); err != nil {
		return utils.UnprocessableEntityResponse(c, err.Error())
	}

	token, admin, err := h.svc.AuthenticateAdmin(body.Email, body.Password)
	if err != nil {
		switch err {
		case service.ErrInvalidCredentials:
			return utils.UnauthorizedResponse(c, "invalid email or password")
		case service.ErrAdminInactive:
			return utils.ForbiddenResponse(c, "admin is inactive")
		default:
			return utils.InternalServerErrorResponse(c, "failed to login")
		}
	}

	return utils.SuccessResponse(c, "login success", map[string]interface{}{
		"token": token,
		"admin": admin,
	})
}

// (optional) helper supaya compile kalau dipakai di routes
// func allowAdminOrSuperadmin(role models.Role) bool {
// 	return role == models.Admins || role == models.Superadmin
// }
//
// func (h *AdminHandler) SomeAdminEndpoint(c echo.Context) error {
// 	role, ok := utils.GetRole(c)
// 	if !ok || !allowAdminOrSuperadmin(role) {
// 		return c.NoContent(http.StatusForbidden)
// 	}
//
// 	// ...existing code...
// 	return c.NoContent(http.StatusOK)
// }
