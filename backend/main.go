package main

import (
	"log"
	"net/http"

	"pickel-backend/routes"
)

func main() {
	router := routes.SetupRoutes()
	log.Println("🚀 Server running on http://localhost:8080")
	http.ListenAndServe(":8080", router)
}
