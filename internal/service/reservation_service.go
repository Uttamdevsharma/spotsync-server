package service

import (
	"errors"
	"spotsync/internal/dto"
	"spotsync/internal/repository"
)

var (
	ErrReservationNotFound = errors.New("reservation not found")
	ErrForbiddenCancel     = errors.New("you can only cancel your own reservations")
	ErrReservationZoneFull = errors.New("parking zone is full")
)

type ReservationService interface {
	CreateReservation(userID uint, req dto.CreateReservationRequest) (*dto.ReservationResponse, error)
	GetMyReservations(userID uint) ([]dto.MyReservationResponse, error)
	CancelReservation(userID uint, userRole string, id uint) error
	GetAllReservations() ([]dto.AdminReservationResponse, error)
}

type reservationService struct {
	reservationRepo repository.ReservationRepository
}

func NewReservationService(reservationRepo repository.ReservationRepository) ReservationService {
	return &reservationService{reservationRepo: reservationRepo}
}

func (s *reservationService) CreateReservation(userID uint, req dto.CreateReservationRequest) (*dto.ReservationResponse, error) {
	res, err := s.reservationRepo.CreateWithLock(userID, req.ZoneID, req.LicensePlate)
	if err != nil {
		if errors.Is(err, repository.ErrZoneFull) {
			return nil, ErrReservationZoneFull
		}
		if errors.Is(err, repository.ErrZoneNotFound) {
			return nil, ErrZoneNotFoundDomain
		}
		return nil, err
	}

	return &dto.ReservationResponse{
		ID:           res.ID,
		UserID:       res.UserID,
		ZoneID:       res.ZoneID,
		LicensePlate: res.LicensePlate,
		Status:       res.Status,
		CreatedAt:    res.CreatedAt,
		UpdatedAt:    res.UpdatedAt,
	}, nil
}

func (s *reservationService) GetMyReservations(userID uint) ([]dto.MyReservationResponse, error) {
	reservations, err := s.reservationRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	res := []dto.MyReservationResponse{}
	for _, r := range reservations {
		res = append(res, dto.MyReservationResponse{
			ID:           r.ID,
			LicensePlate: r.LicensePlate,
			Status:       r.Status,
			Zone: dto.ReservationZoneDetails{
				ID:   r.Zone.ID,
				Name: r.Zone.Name,
				Type: r.Zone.Type,
			},
			CreatedAt:    r.CreatedAt,
		})
	}
	return res, nil
}

func (s *reservationService) CancelReservation(userID uint, userRole string, id uint) error {
	reservation, err := s.reservationRepo.FindByID(id)
	if err != nil {
		return ErrReservationNotFound
	}

	if userRole != "admin" && reservation.UserID != userID {
		return ErrForbiddenCancel
	}

	return s.reservationRepo.UpdateStatus(id, "cancelled")
}

func (s *reservationService) GetAllReservations() ([]dto.AdminReservationResponse, error) {
	reservations, err := s.reservationRepo.FindAll()
	if err != nil {
		return nil, err
	}

	res := []dto.AdminReservationResponse{}
	for _, r := range reservations {
		res = append(res, dto.AdminReservationResponse{
			ID:     r.ID,
			UserID: r.UserID,
			User: dto.ReservationUserDetails{
				ID:    r.User.ID,
				Name:  r.User.Name,
				Email: r.User.Email,
				Role:  r.User.Role,
			},
			ZoneID: r.ZoneID,
			Zone: dto.ReservationZoneDetails{
				ID:   r.Zone.ID,
				Name: r.Zone.Name,
				Type: r.Zone.Type,
			},
			LicensePlate: r.LicensePlate,
			Status:       r.Status,
			CreatedAt:    r.CreatedAt,
			UpdatedAt:    r.UpdatedAt,
		})
	}
	return res, nil
}
