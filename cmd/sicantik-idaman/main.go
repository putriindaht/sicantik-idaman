package main

import (
	"fmt"
	"log"
	"os"
	"sicantik-idaman/configs"
	"sicantik-idaman/internal/db"
	"sicantik-idaman/internal/domain"
	"sicantik-idaman/internal/routes"
	"sicantik-idaman/pkg/databases"
	"sicantik-idaman/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Hello from sicantik-idaman")
	var cfg domain.Helper

	configs.LoadEnv()
	configs.SetHelper(&cfg)

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

	// gin
	r := gin.Default()
	routes.Register(r, cfg)

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(("error run apps sicantik"))
	}

}
