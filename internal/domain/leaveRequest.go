package domain

import (
	"time"

	"github.com/google/uuid"
)

type LeaveRequest struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID       uuid.UUID
	LeaveTypeID  uuid.UUID
	StartDate    time.Time
	EndDate      time.Time
	Reason       string
	Status       LeaveStatus `gorm:"type:varchar(20);default:'pending'"`
	ApprovedByID *uuid.UUID
	ApprovedAt   *time.Time
	NotifyTeam   bool `gorm:"default:false"`
	IsDeleted    bool `gorm:"default:false"`
	CreatedAt    time.Time
	UpdatedAt    time.Time

	User       User            `gorm:"foreignKey:UserID"`
	LeaveType  LeaveType       `gorm:"foreignKey:LeaveTypeID"`
	ApprovedBy *User           `gorm:"foreignKey:ApprovedByID"`
	Reactions  []LeaveReaction `gorm:"foreignKey:LeaveRequestID"`
}
