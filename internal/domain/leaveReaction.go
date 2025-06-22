package domain

import (
	"time"

	"github.com/google/uuid"
)

type LeaveReaction struct {
	ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID         uuid.UUID `json:"user_id"`
	LeaveRequestID uuid.UUID `json:"leave_request_id"`
	Reaction       string    `json:"reaction"`
	IsDeleted      bool      `gorm:"default:false" json:"is_deleted"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	User         User          `gorm:"foreignKey:UserID" json:"-"`
	LeaveRequest *LeaveRequest `gorm:"foreignKey:LeaveRequestID" json:"leave_request,omitempty"`
}
