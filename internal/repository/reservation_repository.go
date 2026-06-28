package repository

import (
	"errors"
	"spotsync/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrZoneFull     = errors.New("parking zone is full")
	ErrZoneNotFound = errors.New("parking zone not found")
)

type ReservationRepository interface {
	CreateWithLock(userID uint, zoneID uint, licensePlate string) (*models.Reservation, error)
	FindByID(id uint) (*models.Reservation, error)
	FindByUserID(userID uint) ([]models.Reservation, error)
	FindAll() ([]models.Reservation, error)
	UpdateStatus(id uint, status string) error
}

type reservationRepository struct {
	db *gorm.DB
}

func NewReservationRepository(db *gorm.DB) ReservationRepository {
	return &reservationRepository{db: db}
}

func (r *reservationRepository) CreateWithLock(userID uint, zoneID uint, licensePlate string) (*models.Reservation, error) {
	var reservation models.Reservation

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var zone models.ParkingZone

		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&zone, zoneID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrZoneNotFound
			}
			return err
		}

		var activeCount int64
		if err := tx.Model(&models.Reservation{}).
			Where("zone_id = ? AND status = 'active'", zoneID).
			Count(&activeCount).Error; err != nil {
			return err
		}

		if int(activeCount) >= zone.TotalCapacity {
			return ErrZoneFull
		}

		reservation = models.Reservation{
			UserID:       userID,
			ZoneID:       zoneID,
			LicensePlate: licensePlate,
			Status:       "active",
		}

		if err := tx.Create(&reservation).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	r.db.Preload("Zone").Preload("User").First(&reservation, reservation.ID)

	return &reservation, nil
}

func (r *reservationRepository) FindByID(id uint) (*models.Reservation, error) {
	var reservation models.Reservation
	err := r.db.Preload("Zone").Preload("User").First(&reservation, id).Error
	if err != nil {
		return nil, err
	}
	return &reservation, nil
}

func (r *reservationRepository) FindByUserID(userID uint) ([]models.Reservation, error) {
	var reservations []models.Reservation
	err := r.db.Preload("Zone").Where("user_id = ?", userID).Order("created_at desc").Find(&reservations).Error
	if err != nil {
		return nil, err
	}
	return reservations, nil
}

func (r *reservationRepository) FindAll() ([]models.Reservation, error) {
	var reservations []models.Reservation
	err := r.db.Preload("Zone").Preload("User").Order("created_at desc").Find(&reservations).Error
	if err != nil {
		return nil, err
	}
	return reservations, nil
}

func (r *reservationRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&models.Reservation{}).Where("id = ?", id).Update("status", status).Error
}
