package main

// @title Darul Abror API
// @version 1.0
// @description Darul Abror backend API (public + admin).
// @BasePath /
// @schemes https http
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer <token>"

import (
	"context"
	"darulabror/api/routes"
	"darulabror/config"
	"darulabror/internal/handler"
	"darulabror/internal/repository"
	"darulabror/internal/service"
	"log"
	"os"
	"reflect"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
)

// CustomValidator enables c.Validate(...) in handlers.
type CustomValidator struct {
	v *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.v.Struct(i)
}

func main() {
	ctx := context.Background()

	// ======================
	// Core (Echo + middleware)
	// ======================
	e := echo.New()
	e.HideBanner = true

	e.Use(echomw.RequestID())
	e.Use(echomw.Recover())

	// CORS (frontend origins) - REQUIRED for production
	corsOrigins := strings.TrimSpace(os.Getenv("CORS_ORIGINS"))
	if corsOrigins == "" {
		log.Fatal("CORS_ORIGINS is required (comma-separated), e.g. https://www.darulabror.com,https://admin.darulabror.com")
	}

	originsRaw := strings.Split(corsOrigins, ",")
	allowOrigins := make([]string, 0, len(originsRaw))
	for _, o := range originsRaw {
		o = strings.TrimSpace(o)
		if o != "" {
			allowOrigins = append(allowOrigins, o)
		}
	}
	if len(allowOrigins) == 0 {
		log.Fatal("CORS_ORIGINS is invalid (no usable origins after parsing)")
	}

	e.Use(echomw.CORSWithConfig(echomw.CORSConfig{
		AllowOrigins: allowOrigins,
		AllowMethods: []string{
			echo.GET,
			echo.POST,
			echo.PUT,
			echo.DELETE,
			echo.OPTIONS,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
		},
		AllowCredentials: false,
	}))

	// Request logging (covers endpoints that return c.NoContent too)
	e.Use(echomw.Logger())

	// Validator for c.Validate(...)
	v := validator.New()
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		if name == "" {
			return fld.Name
		}
		return name
	})
	e.Validator = &CustomValidator{v: v}

	// ======================
	// Database
	// ======================
	db := config.ConnectionDb()

	jwtSecret := strings.TrimSpace(os.Getenv("JWT_SECRET"))
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	// ======================
	// GCS (bucket)
	// ======================
	publicBucket := os.Getenv("PUBLIC_BUCKET")

	var gcsClient *storage.Client
	if publicBucket != "" {
		var err error
		gcsClient, err = storage.NewClient(ctx)
		if err != nil {
			log.Fatalf("failed to init gcs client: %v", err)
		}
	}

	// Always inject (repo will return ErrStorageNotConfigured if not configured)
	publicStore := repository.NewGCPStorageRepo(gcsClient, publicBucket, true)

	// ======================
	// Repositories
	// ======================
	articleRepo := repository.NewArticleRepo(db)
	regRepo := repository.NewRegistrationRepo(db)
	contactRepo := repository.NewContactRepository(db)
	adminRepo := repository.NewAdminRepository(db)

	// ======================
	// Services
	// ======================
	articleSvc := service.NewArticleService(articleRepo, publicStore)
	regSvc := service.NewRegistrationService(regRepo)
	contactSvc := service.NewContactService(contactRepo)
	adminSvc := service.NewAdminService(adminRepo, jwtSecret)

	// ======================
	// Handlers
	// ======================
	h := routes.Handlers{
		Article:      handler.NewArticleHandler(articleSvc),
		Registration: handler.NewRegistrationHandler(regSvc),
		Contact:      handler.NewContactHandler(contactSvc),
		Admin:        handler.NewAdminHandler(adminSvc),
	}

	// ======================
	// Routes
	// ======================
	routes.Register(e, h)

	// ======================
	// Start
	// ======================
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := e.Start(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
	_ = ctx
}
