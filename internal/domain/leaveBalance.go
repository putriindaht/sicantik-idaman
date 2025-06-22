package domain

import (
	"time"

	"github.com/google/uuid"
)

type LeaveBalance struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	LeaveTypeID uuid.UUID `json:"leave_type_id"`
	Year        int       `json:"year"`
	IsPaid      bool      `gorm:"default:false" json:"is_paid"`
	TotalDays   int       `json:"total_days"`
	UsedDays    int       `json:"used_days"`
	IsDeleted   bool      `gorm:"default:false" json:"is_deleted"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	User      User      `gorm:"foreignKey:UserID"      json:"-"`
	LeaveType LeaveType `gorm:"foreignKey:LeaveTypeID" json:"leave_type"`
}
