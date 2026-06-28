package middleware

import (
	"net/http"
	"spotsync/internal/config"
	"spotsync/internal/dto"
	"spotsync/internal/service"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

func JWTMiddleware(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
					Success: false,
					Message: "Unauthorized",
					Errors:  "Missing authorization header",
				})
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
					Success: false,
					Message: "Unauthorized",
					Errors:  "Invalid token format",
				})
			}

			tokenString := parts[1]
			claims := &service.JWTCustomClaims{}

			token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
				return []byte(cfg.JWTSecret), nil
			})

			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
					Success: false,
					Message: "Unauthorized",
					Errors:  "Invalid or expired token",
				})
			}

			c.Set("userID", claims.UserID)
			c.Set("role", claims.Role)

			return next(c)
		}
	}
}

func RequireRole(roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userRole, ok := c.Get("role").(string)
			if !ok {
				return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
					Success: false,
					Message: "Unauthorized",
					Errors:  "User role not found in context",
				})
			}

			allowed := false
			for _, r := range roles {
				if userRole == r {
					allowed = true
					break
				}
			}

			if !allowed {
				return c.JSON(http.StatusForbidden, dto.ErrorResponse{
					Success: false,
					Message: "Forbidden",
					Errors:  "Insufficient permissions",
				})
			}

			return next(c)
		}
	}
}
