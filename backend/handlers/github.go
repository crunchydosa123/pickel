package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"pickel-backend/utils"
)

type GithubPushPayload struct {
	Ref        string `json:"ref"`
	Repository struct {
		CloneURL string `json:"clone_url"`
	} `json:"repository"`
	After string `json:"after"`
}

func GithubWebhook(w http.ResponseWriter, r *http.Request) {
	var payload GithubPushPayload
	json.NewDecoder(r.Body).Decode(&payload)

	if filepath.Base(payload.Ref) != "main" {
		w.WriteHeader(http.StatusOK)
		return
	}

	modelID := "sample"
	go func() {
		if err := utils.TriggerDeploy(payload.Repository.CloneURL, payload.After, modelID); err != nil {
			fmt.Println("Deployment failed:", err)
		}
	}()
}
