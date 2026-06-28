package models

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Email     string    `gorm:"unique;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"`
	Role      string    `gorm:"type:varchar(20);default:'driver';not null" json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ParkingZone struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name          string    `gorm:"not null" json:"name"`
	Type          string    `gorm:"type:varchar(50);not null" json:"type"`
	TotalCapacity int       `gorm:"not null" json:"total_capacity"`
	PricePerHour  float64   `gorm:"type:decimal(10,2);not null" json:"price_per_hour"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Reservation struct {
	ID           uint        `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       uint        `gorm:"not null" json:"user_id"`
	User         User        `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"user,omitempty"`
	ZoneID       uint        `gorm:"not null" json:"zone_id"`
	Zone         ParkingZone `gorm:"foreignKey:ZoneID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"zone,omitempty"`
	LicensePlate string      `gorm:"type:varchar(15);not null" json:"license_plate"`
	Status       string      `gorm:"type:varchar(20);default:'active';not null" json:"status"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}
