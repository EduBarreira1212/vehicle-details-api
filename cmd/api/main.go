package main

import (
	"log"

	"github.com/EduBarreira1212/vehicle-details-api/internal/config"
	apihttp "github.com/EduBarreira1212/vehicle-details-api/internal/http"
)

func main() {
	cfg := config.LoadConfig()

	router := apihttp.BuildRouter()

	addr := ":" + cfg.Port
	log.Printf("listening on %s", addr)

	if err := router.Run(addr); err != nil {
		log.Fatal(err)
	}
}
