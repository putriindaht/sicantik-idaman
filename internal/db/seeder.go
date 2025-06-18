package db

import (
	"log"
	"sicantik-idaman/internal/domain"
	"sicantik-idaman/pkg/databases"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func hashPassword(password string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}
	return string(hashed)
}

func SeedTeamAndUser() {
	db := databases.DB

	// Step 1: Seed Teams
	teamNames := []string{"Engineering", "Human Resource", "Finance"}

	teamMap := make(map[string]uuid.UUID)

	for _, name := range teamNames {
		team := domain.Team{
			ID:   uuid.New(),
			Name: name,
		}

		// Avoid duplicate team
		var existing domain.Team
		if err := db.Where("name = ?", name).First(&existing).Error; err == gorm.ErrRecordNotFound {
			db.Create(&team)
			teamMap[name] = team.ID
		} else {
			teamMap[name] = existing.ID
		}
	}

	engID := teamMap["Engineering"]
	hrID := teamMap["Human Resource"]
	financeID := teamMap["Finance"]

	// Step 2: Seed Users
	users := []domain.User{
		{
			ID:       uuid.New(),
			Name:     "Eva Engineer",
			Email:    "eva.engineer@example.com",
			Password: hashPassword("password123"),
			Role:     "employee",
			Title:    "Backend Engineer",
			TeamID:   &engID,
		},
		{
			ID:       uuid.New(),
			Name:     "Mike Manager",
			Email:    "mike.manager@example.com",
			Password: hashPassword("password123"),
			Role:     "manager",
			Title:    "Engineering Manager",
			TeamID:   &engID,
		},
		{
			ID:       uuid.New(),
			Name:     "Helen HR",
			Email:    "helen.hr@example.com",
			Password: hashPassword("password123"),
			Role:     "hr",
			Title:    "HR Specialist",
			TeamID:   &hrID,
		},
		{
			ID:       uuid.New(),
			Name:     "Frank Finance",
			Email:    "frank.finance@example.com",
			Password: hashPassword("password123"),
			Role:     "employee",
			Title:    "Finance Analyst",
			TeamID:   &financeID,
		},
		{
			ID:       uuid.New(),
			Name:     "Diana Director",
			Email:    "diana.director@example.com",
			Password: hashPassword("password123"),
			Role:     "director",
			Title:    "Director of Operations",
			TeamID:   &hrID,
		},
	}

	for _, user := range users {
		var existing domain.User
		if err := db.Where("email = ?", user.Email).First(&existing).Error; err == gorm.ErrRecordNotFound {
			db.Create(&user)
		}
	}

	log.Println("Seeder Team and User complete ✅")
}

func SeedLeaveTypes() map[string]uuid.UUID {
	db := databases.DB

	leaveTypes := []domain.LeaveType{
		{
			ID:               uuid.New(),
			Name:             "Cuti Tahunan",
			Description:      "Cuti tahunan yang diberikan setiap tahun",
			DefaultDays:      12,
			IsPaid:           true,
			RequiresApproval: true,
			RequiresDocument: false,
		},
		{
			ID:               uuid.New(),
			Name:             "Cuti Sakit",
			Description:      "Cuti karena sakit dengan surat keterangan dokter",
			DefaultDays:      14,
			IsPaid:           true,
			RequiresApproval: true,
			RequiresDocument: true,
		},
		{
			ID:               uuid.New(),
			Name:             "Cuti Melahirkan",
			Description:      "Cuti untuk melahirkan",
			DefaultDays:      90,
			IsPaid:           true,
			RequiresApproval: true,
			RequiresDocument: true,
		},
		{
			ID:               uuid.New(),
			Name:             "Cuti Menikah",
			Description:      "Cuti untuk pernikahan",
			DefaultDays:      3,
			IsPaid:           true,
			RequiresApproval: true,
			RequiresDocument: false,
		},
	}

	leaveTypeMap := make(map[string]uuid.UUID)

	for _, lt := range leaveTypes {
		var existing domain.LeaveType
		if err := db.Where("name = ?", lt.Name).First(&existing).Error; err != nil {
			db.Create(&lt)
			leaveTypeMap[lt.Name] = lt.ID
		} else {
			leaveTypeMap[lt.Name] = existing.ID
		}
	}

	log.Println("✅ Seeded Leave Types")
	return leaveTypeMap
}

func SeedLeaveBalances(leaveTypeMap map[string]uuid.UUID) {
	db := databases.DB

	// get all users
	var users []domain.User
	if err := db.Find(&users).Error; err != nil {
		log.Fatalf("Failed to fetch users: %v", err)
	}

	currentYear := time.Now().Year()

	for _, user := range users {
		for _, leaveTypeID := range leaveTypeMap {
			var existing domain.LeaveBalance
			err := db.Where("user_id = ? AND leave_type_id = ? AND year = ?", user.ID, leaveTypeID, currentYear).First(&existing).Error
			if err == nil {
				continue // already exists
			}

			// Get default days from type
			var leaveType domain.LeaveType
			if err := db.First(&leaveType, "id = ?", leaveTypeID).Error; err != nil {
				continue
			}

			lb := domain.LeaveBalance{
				ID:          uuid.New(),
				UserID:      user.ID,
				LeaveTypeID: leaveTypeID,
				Year:        currentYear,
				IsPaid:      leaveType.IsPaid,
				TotalDays:   leaveType.DefaultDays,
				UsedDays:    0,
			}
			db.Create(&lb)
		}
	}

	log.Println("✅ Seeded Leave Balances")
}

func SeedLeaveData() {
	leaveTypeMap := SeedLeaveTypes()
	SeedLeaveBalances(leaveTypeMap)
}

func SeedData() {
	SeedTeamAndUser()
	SeedLeaveData()
}
