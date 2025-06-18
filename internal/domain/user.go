package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name      string
	Email     string `gorm:"unique"`
	Password  string
	Role      UserRole `gorm:"default:'employee'"`
	Title     string
	IsDeleted bool `gorm:"default:false"`
	CreatedAt time.Time
	UpdatedAt time.Time

	TeamID        *uuid.UUID
	LeaveBalances []LeaveBalance  `gorm:"foreignKey:UserID"`
	LeaveRequests []LeaveRequest  `gorm:"foreignKey:UserID"`
	Reactions     []LeaveReaction `gorm:"foreignKey:UserID"`
}
