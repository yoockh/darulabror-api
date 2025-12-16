package routes

import (
	"darulabror/api/middleware"
	"darulabror/internal/handler"
	"darulabror/internal/models"

	"github.com/labstack/echo/v4"
)

type Handlers struct {
	Article      *handler.ArticleHandler
	Registration *handler.RegistrationHandler
	Contact      *handler.ContactHandler
	Admin        *handler.AdminHandler
}

func Register(e *echo.Echo, h Handlers) {
	// ======================
	// Public routes
	// ======================
	e.GET("/articles", h.Article.ListPublished)
	e.GET("/articles/:id", h.Article.GetPublishedByID)

	e.POST("/registrations", h.Registration.Create)
	e.POST("/contacts", h.Contact.Create)

	// ======================
	// Admin routes (/admin)
	// ======================
	admin := e.Group("/admin", middleware.JWTAuth(), middleware.RequireRole(models.Admins, models.Superadmin))

	admin.GET("/profile", h.Admin.Profile)

	// manage articles
	admin.GET("/articles", h.Article.AdminListAll)
	admin.POST("/articles", h.Article.AdminCreate)
	admin.PUT("/articles/:id", h.Article.AdminUpdate)
	admin.DELETE("/articles/:id", h.Article.AdminDelete)

	// manage registrations
	admin.GET("/registrations", h.Registration.AdminList)
	admin.GET("/registrations/:id", h.Registration.AdminGetByID)
	admin.DELETE("/registrations/:id", h.Registration.AdminDelete)

	// manage contacts
	admin.GET("/contacts", h.Contact.AdminList)
	admin.GET("/contacts/:id", h.Contact.AdminGetByID)
	admin.PUT("/contacts/:id", h.Contact.AdminUpdate)
	admin.DELETE("/contacts/:id", h.Contact.AdminDelete)

	// ======================
	// Superadmin-only routes
	// ======================
	super := admin.Group("", middleware.RequireRole(models.Superadmin))
	super.POST("/admins", h.Admin.Create)
	super.GET("/admins", h.Admin.List)
	super.PUT("/admins/:id", h.Admin.Update)
	super.DELETE("/admins/:id", h.Admin.Delete)
}
