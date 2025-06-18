package domain

import (
	"time"

	"github.com/google/uuid"
)

type LeaveReaction struct {
	ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID         uuid.UUID
	LeaveRequestID uuid.UUID
	Reaction       string
	IsDeleted      bool `gorm:"default:false"`
	CreatedAt      time.Time
	UpdatedAt      time.Time

	User         User         `gorm:"foreignKey:UserID"`
	LeaveRequest LeaveRequest `gorm:"foreignKey:LeaveRequestID"`
}
