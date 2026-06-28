package routes

import (
	"spotsync/internal/config"
	"spotsync/internal/handler"
	"spotsync/internal/middleware"

	"github.com/labstack/echo/v4"
)

func Register(e *echo.Group, authHandler *handler.AuthHandler, zoneHandler *handler.ZoneHandler, resHandler *handler.ReservationHandler, cfg *config.Config) {
	v1 := e

	v1.POST("/auth/register", authHandler.Register)
	v1.POST("/auth/login", authHandler.Login)

	v1.GET("/zones", zoneHandler.GetAllZones)
	v1.GET("/zones/:id", zoneHandler.GetZoneByID)

	protected := v1.Group("")
	protected.Use(middleware.JWTMiddleware(cfg))

	protected.POST("/zones", zoneHandler.CreateZone, middleware.RequireRole("admin"))

	protected.POST("/reservations", resHandler.CreateReservation, middleware.RequireRole("driver", "admin"))
	protected.GET("/reservations/my-reservations", resHandler.GetMyReservations, middleware.RequireRole("driver", "admin"))
	protected.DELETE("/reservations/:id", resHandler.CancelReservation, middleware.RequireRole("driver", "admin"))

	protected.GET("/reservations", resHandler.GetAllReservations, middleware.RequireRole("admin"))
}
