package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pickel-backend/middleware"
	"pickel-backend/utils"
	"time"

	"github.com/gorilla/mux"
)

type CreateModelRequest struct {
	Name string `json:"name"`
}

type DeployModelRequest struct {
	FileName string `json:"file_name"`
	ModelId  string `json:"model_id"`
	Content  []byte `json:"content"`
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

func DeployModel(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetUserFromContext(r)
	userId := claims.UserID
	modelID := r.FormValue("modelId")

	err := r.ParseMultipartForm(20 << 20)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "failed to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("modelFile")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "file missing: "+err.Error(), http.StatusBadRequest)
		return
	}

	defer file.Close()

	s3Key := fmt.Sprintf("%s/%s", userId, header.Filename)

	s3URL, err := utils.UploadFileToS3(r.Context(), file, s3Key)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "S3 upload failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	lambdaArn, err := utils.DeployToLambda(header.Filename, s3Key)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Lambda deployment failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	apiURL, err := utils.CreateAPIGatewayForLambda(r.Context(), lambdaArn)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "API Gateway deployment failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("s3URL", s3URL)
	fmt.Println("lambda endpoint", apiURL)

	db := utils.GetDB()
	_, err = db.Exec(context.Background(),
		"INSERT INTO modelcode (s3URL, model_id, lambda_arn, api_url) VALUES ($1, $2, $3, $4)",
		s3URL, modelID, lambdaArn, apiURL,
	)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "DB insert failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":     "Model deployed successfully",
		"modelId":     modelID,
		"apiEndpoint": apiURL,
		"s3URL":       s3URL,
		"fileName":    header.Filename,
		"createdBy":   userId,
	})

}

func GetSingleModel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	db := utils.GetDB()

	type Model struct {
		ID             string    `json:"id"`
		Name           string    `json:"name"`
		UserID         string    `json:"user_id"`
		CreatedAt      time.Time `json:"created_at"`
		DeploymentType string    `json:"deployment_type"`
	}

	var model Model

	err := db.QueryRow(
		context.Background(),
		"SELECT id, name, created_by, created_at, deployment_type FROM models WHERE id = $1",
		id,
	).Scan(&model.ID, &model.Name, &model.UserID, &model.CreatedAt, &model.DeploymentType)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Model not found: "+err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(model)
}

// get public url
func GetModelURL(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"url": "https://example.com/model.pkl"})
}
