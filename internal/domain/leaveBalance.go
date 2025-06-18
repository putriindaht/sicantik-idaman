package domain

import (
	"time"

	"github.com/google/uuid"
)

type LeaveBalance struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID      uuid.UUID
	LeaveTypeID uuid.UUID
	Year        int
	IsPaid      bool `gorm:"default:false"`
	TotalDays   int
	UsedDays    int
	IsDeleted   bool `gorm:"default:false"`
	CreatedAt   time.Time
	UpdatedAt   time.Time

	User      User      `gorm:"foreignKey:UserID"`
	LeaveType LeaveType `gorm:"foreignKey:LeaveTypeID"`
}
