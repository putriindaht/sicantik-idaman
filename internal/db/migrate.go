package db

import (
	"sicantik-idaman/internal/domain"
	"sicantik-idaman/pkg/databases"
)

func Migrate() error {
	return databases.DB.AutoMigrate(
		&domain.Team{},
		&domain.User{},
		&domain.LeaveRequest{},
		&domain.LeaveBalance{},
		&domain.LeaveType{},
		&domain.LeaveReaction{},
	)
}
