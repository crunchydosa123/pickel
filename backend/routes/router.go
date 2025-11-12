package routes

import (
	"fmt"
	"net/http"
	"pickel-backend/handlers"
	"pickel-backend/middleware"

	"github.com/gorilla/mux"
)

func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/auth/signup", handlers.Signup).Methods("POST")
	r.HandleFunc("/auth/login", handlers.Login).Methods("POST")

	modelSubrouter := r.PathPrefix("/model").Subrouter()
	modelSubrouter.Use(middleware.JWTAuth)

	modelSubrouter.HandleFunc("/create", handlers.CreateModel).Methods("POST")
	modelSubrouter.HandleFunc("/add", handlers.AddFileToModel).Methods("POST")
	modelSubrouter.HandleFunc("/deploy", handlers.DeployModel).Methods("POST")
	modelSubrouter.HandleFunc("/url", handlers.GetModelURL).Methods("GET")
	modelSubrouter.HandleFunc("/", handlers.GetModelByUser).Methods("GET")

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("Incoming request:", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	})

	return r
}
