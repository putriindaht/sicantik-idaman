package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name      string    `json:"name"`
	Email     string    `gorm:"unique" json:"email"`
	Password  string    `json:"-"`
	Role      UserRole  `gorm:"default:'employee'" json:"role"`
	Title     string    `json:"title"`
	IsDeleted bool      `gorm:"default:false" json:"is_deleted"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	TeamID        *uuid.UUID      `json:"team_id"`
	LeaveBalances []LeaveBalance  `gorm:"foreignKey:UserID" json:"leave_balances"`
	LeaveRequests []LeaveRequest  `gorm:"foreignKey:UserID" json:"leave_requests"`
	Reactions     []LeaveReaction `gorm:"foreignKey:UserID" json:"reactions"`
}
