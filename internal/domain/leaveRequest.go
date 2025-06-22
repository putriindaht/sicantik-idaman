package domain

import (
	"time"

	"github.com/google/uuid"
)

type LeaveRequest struct {
	ID           uuid.UUID   `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID       uuid.UUID   `json:"user_id"`
	LeaveTypeID  uuid.UUID   `json:"leave_type_id"`
	StartDate    time.Time   `json:"start_date"`
	EndDate      time.Time   `json:"end_date"`
	Reason       string      `json:"reason"`
	Status       LeaveStatus `gorm:"type:varchar(20);default:'pending'" json:"status"`
	ApprovedByID *uuid.UUID  `json:"approve_by_id"`
	ApprovedAt   *time.Time  `json:"approve_at"`
	RejectedNote *string     `gorm:"type:text" json:"rejected_note"`
	NotifyTeam   bool        `gorm:"default:false" json:"notify_team"`
	IsDeleted    bool        `gorm:"default:false" json:"is_deleted"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`

	User       User            `gorm:"foreignKey:UserID" json:"-"`
	LeaveType  *LeaveType      `gorm:"foreignKey:LeaveTypeID" json:"leave_type,omitempty"`
	ApprovedBy *User           `gorm:"foreignKey:ApprovedByID" json:"-"`
	Reactions  []LeaveReaction `gorm:"foreignKey:LeaveRequestID" json:"reactions,omitempty"`
}
