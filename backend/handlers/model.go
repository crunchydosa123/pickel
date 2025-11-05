package handlers

import (
	"encoding/json"
	"net/http"
	"pickel-backend/utils"
)

// create a new model
func CreateModel(w http.ResponseWriter, r *http.Request) {
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

// add a file to a model (h5, pickel)
func AddFileToModel(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"message": "Model Added"})
}

// allow the model to get requests from public
func DeployModel(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"message": "Model deployed"})
}

// get public url
func GetModelURL(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"url": "https://example.com/model.pkl"})
}
