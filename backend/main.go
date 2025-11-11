package main

import (
	"fmt"
	"net/http"
	"pickel-backend/routes"
)

func main() {
	r := routes.SetupRoutes()

	fmt.Println("✅ Server running on http://localhost:8080")
	err := http.ListenAndServe(":8080", routes.EnableCORS(r))
	if err != nil {
		fmt.Println("Server error:", err)
	}
}
