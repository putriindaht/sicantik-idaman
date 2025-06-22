package domain

import (
	"time"

	"github.com/google/uuid"
)

type LeaveType struct {
	ID               uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	DefaultDays      int       `json:"default_days"`
	IsPaid           bool      `gorm:"default:false" json:"is_paid"`
	RequiresApproval bool      `gorm:"default:false" json:"requires_approval"`
	RequiresDocument bool      `gorm:"default:false" json:"requires_document"`
	IsDeleted        bool      `gorm:"default:false" json:"is_deleted"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	LeaveBalances []LeaveBalance `gorm:"foreignKey:LeaveTypeID" json:"-"`
	LeaveRequests []LeaveRequest `gorm:"foreignKey:LeaveTypeID" json:"-"`
}
