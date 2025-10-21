package main

import (
	"log"

	"github.com/EduBarreira1212/vehicle-details-api/internal/config"
	apihttp "github.com/EduBarreira1212/vehicle-details-api/internal/http"
	"github.com/EduBarreira1212/vehicle-details-api/internal/models"
)

func main() {
	cfg := config.LoadConfig()
	config.ConnectDatabase(cfg.DB_URL)

	config.DB.AutoMigrate(&models.User{}, &models.History{})

	router := apihttp.BuildRouter()

	addr := ":" + cfg.Port
	log.Printf("listening on %s", addr)

	if err := router.Run(addr); err != nil {
		log.Fatal(err)
	}
}
