package handler

import (
	"errors"
	"net/http"
	"strconv"
	"spotsync/internal/dto"
	"spotsync/internal/service"

	"github.com/labstack/echo/v4"
)

type ZoneHandler struct {
	zoneService service.ZoneService
}

func NewZoneHandler(zoneService service.ZoneService) *ZoneHandler {
	return &ZoneHandler{zoneService: zoneService}
}

func (h *ZoneHandler) CreateZone(c echo.Context) error {
	req := new(dto.CreateZoneRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Invalid request payload",
			Errors:  err.Error(),
		})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Validation failed",
			Errors:  err.Error(),
		})
	}

	res, err := h.zoneService.CreateZone(*req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "Failed to create parking zone",
			Errors:  err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, dto.SuccessResponse{
		Success: true,
		Message: "Parking zone created successfully",
		Data:    res,
	})
}

func (h *ZoneHandler) GetAllZones(c echo.Context) error {
	res, err := h.zoneService.GetAllZones()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "Failed to retrieve parking zones",
			Errors:  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "Parking zones retrieved successfully",
		Data:    res,
	})
}

func (h *ZoneHandler) GetZoneByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Invalid zone ID",
			Errors:  err.Error(),
		})
	}

	res, err := h.zoneService.GetZoneByID(uint(id))
	if err != nil {
		if errors.Is(err, service.ErrZoneNotFoundDomain) {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Success: false,
				Message: "Parking zone not found",
				Errors:  err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "Failed to retrieve parking zone",
			Errors:  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "Parking zone retrieved successfully",
		Data:    res,
	})
}
