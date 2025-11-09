package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"pickel-backend/middleware"
	"pickel-backend/utils"
)

type CreateModelRequest struct {
	Name       string `json:"name"`
	Created_By string `json:"created_by"`
}

// create a new model
func CreateModel(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID := claims.UserID
	fmt.Fprintf(w, "Model created by user: %s", userID)

	var req CreateModelRequest
	var db = utils.GetDB()
	_, err := db.Exec(context.Background(),
		"INSERT INTO models values ($1, $2)",
		req.Name, req.Created_By,
	)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Model created"})

}

// add a file to a model (h5, pickel)
func AddFileToModel(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Missing file: "+err.Error(), http.StatusBadRequest)
		return
	}

	defer file.Close()

	url, err := utils.UploadFileToS3(r.Context(), file, handler.Filename)

	if err != nil {
		http.Error(w, "S3 upload failed: "+err.Error(), http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "File uploaded successfully",
		"url":     url,
	})
}

// allow the model to get requests from public
func DeployModel(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"message": "Model deployed"})
}

// get public url
func GetModelURL(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"url": "https://example.com/model.pkl"})
}
