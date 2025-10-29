package main

import (
	"log"

	"github.com/EduBarreira1212/vehicle-details-api/internal/config"
	apihttp "github.com/EduBarreira1212/vehicle-details-api/internal/http"
	"github.com/EduBarreira1212/vehicle-details-api/internal/migrations"
	"github.com/EduBarreira1212/vehicle-details-api/internal/models"
)

func main() {
	cfg := config.LoadConfig()
	config.ConnectDatabase(cfg.DATABASE_URL)

	m := migrations.New(config.DB)
	if err := m.Migrate(); err != nil {
		log.Fatalf("migrations failed: %v", err)
	}
	log.Println("migrations applied")

	config.DB.AutoMigrate(&models.User{}, &models.History{})

	router := apihttp.BuildRouter()

	addr := ":" + cfg.Port
	log.Printf("listening on %s", addr)

	if err := router.Run(addr); err != nil {
		log.Fatal(err)
	}
}
