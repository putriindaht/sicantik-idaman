package domain

import (
	"time"

	"github.com/google/uuid"
)

type Team struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name      string    `json:"name"`
	IsDeleted bool      `gorm:"default:false" json:"is_deleted"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Users     []User    `gorm:"foreignKey:TeamID" json:"-"`
}
