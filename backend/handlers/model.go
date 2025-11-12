package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pickel-backend/middleware"
	"pickel-backend/utils"
)

type CreateModelRequest struct {
	Name string `json:"name"`
}

func CreateModel(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetUserFromContext(r)

	userID := claims.UserID

	var req CreateModelRequest

	bodyBytes, _ := io.ReadAll(r.Body)
	fmt.Println("Raw body:", string(bodyBytes))
	if err := json.Unmarshal(bodyBytes, &req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var db = utils.GetDB()
	_, err := db.Exec(context.Background(),
		"INSERT INTO models (name, created_by) values ($1, $2)",
		req.Name, userID,
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

func GetModelByUser(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetUserFromContext(r)
	userId := claims.UserID

	utils.ConnectSupabase()
	var db = utils.GetDB()

	rows, err := db.Query(context.Background(),
		"SELECT id, name FROM models WHERE created_by = $1",
		userId)

	if err != nil {
		http.Error(w, fmt.Sprintf("error fetching models for user: %v", err), http.StatusInternalServerError)
	}

	defer rows.Close()

	var models []map[string]interface{}

	for rows.Next() {
		var id, name string
		if err := rows.Scan(&id, &name); err != nil {
			http.Error(w, fmt.Sprintf("error scanning row: %v", err), http.StatusInternalServerError)
			return
		}
		models = append(models, map[string]interface{}{
			"id":   id,
			"name": name,
		})
	}

	if err = rows.Err(); err != nil {
		http.Error(w, fmt.Sprintf("error scanning row: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"userId": userId,
		"models": models,
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
