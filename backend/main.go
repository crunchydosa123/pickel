package main

import (
	"fmt"
	"net/http"
	"os"
	"pickel-backend/routes"
)

func main() {
	r := routes.SetupRoutes()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("Server running on http://localhost:" + port)
	err := http.ListenAndServe(":"+port, routes.EnableCORS(r))
	if err != nil {
		fmt.Println("Server error:", err)
	}

}
