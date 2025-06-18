package databases

import (
	"log"
	"sicantik-idaman/pkg/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connection(dsn string) {
	gormDb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Gagal Koneksi ke database server : ", err)
	}
	DB = gormDb

	// Automatically create the uuid-ossp extension
	DB.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)

	log := logger.Log

	log.Info("Success connection to postgresql database")
}
