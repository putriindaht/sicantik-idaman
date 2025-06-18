package main

import (
	"fmt"
	"log"
	"os"
	"sicantik-idaman/configs"
	"sicantik-idaman/internal/db"
	"sicantik-idaman/pkg/databases"
	"sicantik-idaman/pkg/logger"
)

func main() {
	fmt.Println("Hello from sicantik-idaman")

	configs.LoadEnv()

	logger.InitLoggerZap()

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	databases.Connection(dsn)

	// run migration
	if err := db.Migrate(); err != nil {
		log.Fatal("Migration failed:", err)
	}

	// run seeder
	db.SeedData()

}
