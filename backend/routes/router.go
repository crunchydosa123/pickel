package routes

import (
	"pickel-backend/handlers"
	"pickel-backend/middleware"

	"github.com/gorilla/mux"
)

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

	return r
}
