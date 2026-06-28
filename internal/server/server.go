package server

import (
	"net/http"
	"spotsync/internal/config"
	"spotsync/internal/handler"
	"spotsync/internal/middleware"
	"spotsync/internal/repository"
	"spotsync/internal/routes"
	"spotsync/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

func Start(db *gorm.DB, cfg *config.Config) {
	userRepo := repository.NewUserRepository(db)
	zoneRepo := repository.NewZoneRepository(db)
	resRepo := repository.NewReservationRepository(db)

	authService := service.NewAuthService(userRepo, cfg)
	zoneService := service.NewZoneService(zoneRepo)
	resService := service.NewReservationService(resRepo)

	authHandler := handler.NewAuthHandler(authService)
	zoneHandler := handler.NewZoneHandler(zoneService)
	resHandler := handler.NewReservationHandler(resService)

	e := echo.New()

	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))

	e.Validator = &middleware.CustomValidator{Validator: validator.New()}

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": true,
			"message": "SpotSync API is running...",
		})
	})

	v1 := e.Group("/api/v1")
	routes.Register(v1, authHandler, zoneHandler, resHandler, cfg)

	if err := e.Start(":" + cfg.ServerPort); err != nil && err != http.ErrServerClosed {
		e.Logger.Fatal(err)
	}
}
