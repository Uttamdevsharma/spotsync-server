package handler

import (
	"errors"
	"net/http"
	"strconv"
	"spotsync/internal/dto"
	"spotsync/internal/service"

	"github.com/labstack/echo/v4"
)

type ReservationHandler struct {
	reservationService service.ReservationService
}

func NewReservationHandler(reservationService service.ReservationService) *ReservationHandler {
	return &ReservationHandler{reservationService: reservationService}
}

func (h *ReservationHandler) CreateReservation(c echo.Context) error {
	userID, ok := c.Get("userID").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Success: false,
			Message: "Unauthorized",
			Errors:  "User context not found",
		})
	}

	req := new(dto.CreateReservationRequest)
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

	res, err := h.reservationService.CreateReservation(userID, *req)
	if err != nil {
		if errors.Is(err, service.ErrReservationZoneFull) {
			return c.JSON(http.StatusConflict, dto.ErrorResponse{
				Success: false,
				Message: "Reservation failed",
				Errors:  err.Error(),
			})
		}
		if errors.Is(err, service.ErrZoneNotFoundDomain) {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Success: false,
				Message: "Parking zone not found",
				Errors:  err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "Failed to confirm reservation",
			Errors:  err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, dto.SuccessResponse{
		Success: true,
		Message: "Reservation confirmed successfully",
		Data:    res,
	})
}

func (h *ReservationHandler) GetMyReservations(c echo.Context) error {
	userID, ok := c.Get("userID").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Success: false,
			Message: "Unauthorized",
			Errors:  "User context not found",
		})
	}

	res, err := h.reservationService.GetMyReservations(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "Failed to retrieve reservations",
			Errors:  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "My reservations retrieved successfully",
		Data:    res,
	})
}

func (h *ReservationHandler) CancelReservation(c echo.Context) error {
	userID, ok := c.Get("userID").(uint)
	role, roleOk := c.Get("role").(string)
	if !ok || !roleOk {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Success: false,
			Message: "Unauthorized",
			Errors:  "User context not found",
		})
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Invalid reservation ID",
			Errors:  err.Error(),
		})
	}

	err = h.reservationService.CancelReservation(userID, role, uint(id))
	if err != nil {
		if errors.Is(err, service.ErrReservationNotFound) {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Success: false,
				Message: "Reservation not found",
				Errors:  err.Error(),
			})
		}
		if errors.Is(err, service.ErrForbiddenCancel) {
			return c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Success: false,
				Message: "Forbidden action",
				Errors:  err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "Failed to cancel reservation",
			Errors:  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "Reservation cancelled successfully",
	})
}

func (h *ReservationHandler) GetAllReservations(c echo.Context) error {
	res, err := h.reservationService.GetAllReservations()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "Failed to retrieve reservations",
			Errors:  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "Reservations retrieved successfully",
		Data:    res,
	})
}
