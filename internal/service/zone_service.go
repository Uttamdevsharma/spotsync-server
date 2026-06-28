package service

import (
	"errors"
	"spotsync/internal/dto"
	"spotsync/internal/models"
	"spotsync/internal/repository"
)

var (
	ErrZoneNotFoundDomain = errors.New("zone not found")
)

type ZoneService interface {
	CreateZone(req dto.CreateZoneRequest) (*dto.ZoneResponse, error)
	GetAllZones() ([]dto.ZoneListResponse, error)
	GetZoneByID(id uint) (*dto.ZoneListResponse, error)
}

type zoneService struct {
	zoneRepo repository.ZoneRepository
}

func NewZoneService(zoneRepo repository.ZoneRepository) ZoneService {
	return &zoneService{zoneRepo: zoneRepo}
}

func (s *zoneService) CreateZone(req dto.CreateZoneRequest) (*dto.ZoneResponse, error) {
	zone := &models.ParkingZone{
		Name:          req.Name,
		Type:          req.Type,
		TotalCapacity: req.TotalCapacity,
		PricePerHour:  req.PricePerHour,
	}

	if err := s.zoneRepo.Create(zone); err != nil {
		return nil, err
	}

	return &dto.ZoneResponse{
		ID:            zone.ID,
		Name:          zone.Name,
		Type:          zone.Type,
		TotalCapacity: zone.TotalCapacity,
		PricePerHour:  zone.PricePerHour,
		CreatedAt:     zone.CreatedAt,
		UpdatedAt:     zone.UpdatedAt,
	}, nil
}

func (s *zoneService) GetAllZones() ([]dto.ZoneListResponse, error) {
	zones, err := s.zoneRepo.FindAll()
	if err != nil {
		return nil, err
	}

	res := []dto.ZoneListResponse{}
	for _, zone := range zones {
		activeCount, err := s.zoneRepo.GetActiveCount(zone.ID)
		if err != nil {
			return nil, err
		}

		availableSpots := zone.TotalCapacity - activeCount
		if availableSpots < 0 {
			availableSpots = 0
		}

		res = append(res, dto.ZoneListResponse{
			ID:             zone.ID,
			Name:           zone.Name,
			Type:           zone.Type,
			TotalCapacity:  zone.TotalCapacity,
			AvailableSpots: availableSpots,
			PricePerHour:   zone.PricePerHour,
			CreatedAt:      zone.CreatedAt,
		})
	}
	return res, nil
}

func (s *zoneService) GetZoneByID(id uint) (*dto.ZoneListResponse, error) {
	zone, err := s.zoneRepo.FindByID(id)
	if err != nil {
		return nil, ErrZoneNotFoundDomain
	}

	activeCount, err := s.zoneRepo.GetActiveCount(zone.ID)
	if err != nil {
		return nil, err
	}

	availableSpots := zone.TotalCapacity - activeCount
	if availableSpots < 0 {
		availableSpots = 0
	}

	return &dto.ZoneListResponse{
		ID:             zone.ID,
		Name:           zone.Name,
		Type:           zone.Type,
		TotalCapacity:  zone.TotalCapacity,
		AvailableSpots: availableSpots,
		PricePerHour:   zone.PricePerHour,
		CreatedAt:      zone.CreatedAt,
	}, nil
}
