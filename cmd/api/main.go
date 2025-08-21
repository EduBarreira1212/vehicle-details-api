package main

import (
	"log"

	"github.com/EduBarreira1212/vehicle-details-api/internal/http"
)

func main() {
	router := http.BuildRouter()

	log.Println("listening on port 8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
