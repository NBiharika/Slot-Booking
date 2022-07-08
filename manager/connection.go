package manager

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var dbClient *gorm.DB

func InitilizeConnection() {
	dbconn := "host=localhost user=niharika dbname=postgres port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dbconn), &gorm.Config{})
	if err != nil {
		log.Println("Error", err)
	}
	dbClient = db
}
