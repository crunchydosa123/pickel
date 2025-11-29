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
		port = "8080" // fallback for local dev
	}

	// log port binding
	fmt.Println("Server running on http://0.0.0.0:" + port)

	// start server immediately
	if err := http.ListenAndServe("0.0.0.0:"+port, routes.EnableCORS(r)); err != nil {
		fmt.Println("Server error:", err)
	}
}
