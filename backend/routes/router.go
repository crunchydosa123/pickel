package routes

import (
	"pickel-backend/handlers"

	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/model/create", handlers.CreateModel).Methods("POST")
	r.HandleFunc("/model/add", handlers.AddFileToModel).Methods("POST")
	r.HandleFunc("/model/deploy", handlers.DeployModel).Methods("POST")
	r.HandleFunc("/model/url", handlers.GetModelURL).Methods("GET")

	return r
}
