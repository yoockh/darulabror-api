package middleware

import (
	"darulabror/internal/models"
	"darulabror/internal/utils"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type Claims struct {
	AdminID uint        `json:"admin_id"`
	Role    models.Role `json:"role"`
	jwt.RegisteredClaims
}

func JWTAuth() echo.MiddlewareFunc {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		logrus.Warn("JWT_SECRET is empty (JWT auth will fail)")
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			h := c.Request().Header.Get("Authorization")
			if h == "" || !strings.HasPrefix(h, "Bearer ") {
				return utils.UnauthorizedResponse(c, "missing bearer token")
			}
			tokenStr := strings.TrimPrefix(h, "Bearer ")

			token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})
			if err != nil || !token.Valid {
				logrus.WithError(err).Warn("invalid jwt token")
				return utils.UnauthorizedResponse(c, "invalid token")
			}

			claims, ok := token.Claims.(*Claims)
			if !ok || claims.AdminID == 0 || claims.Role == "" {
				return utils.UnauthorizedResponse(c, "invalid token claims")
			}

			c.Set(utils.CtxAdminIDKey, claims.AdminID)
			c.Set(utils.CtxRoleKey, claims.Role)

			return next(c)
		}
	}
}

func RequireRole(allowed ...models.Role) echo.MiddlewareFunc {
	allowedSet := map[models.Role]bool{}
	for _, r := range allowed {
		allowedSet[r] = true
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			role, ok := utils.GetRole(c)
			if !ok {
				return c.NoContent(http.StatusUnauthorized)
			}
			if !allowedSet[role] {
				logrus.WithField("role", role).Warn("forbidden role")
				return utils.ForbiddenResponse(c, "forbidden")
			}
			return next(c)
		}
	}
}
