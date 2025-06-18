package domain

import (
	"time"

	"github.com/google/uuid"
)

type LeaveType struct {
	ID               uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name             string
	Description      string
	DefaultDays      int
	IsPaid           bool `gorm:"default:false"`
	RequiresApproval bool `gorm:"default:false"`
	RequiresDocument bool `gorm:"default:false"`
	IsDeleted        bool `gorm:"default:false"`
	CreatedAt        time.Time
	UpdatedAt        time.Time

	LeaveBalances []LeaveBalance `gorm:"foreignKey:LeaveTypeID"`
	LeaveRequests []LeaveRequest `gorm:"foreignKey:LeaveTypeID"`
}
