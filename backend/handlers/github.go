package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"pickel-backend/utils"
	"strconv"
	"strings"
)

type GithubPushPayload struct {
	Ref string `json:"ref"`

	Repository struct {
		CloneURL string `json:"clone_url"`
	} `json:"repository"`

	After string `json:"after"`

	Installation struct {
		ID int64 `json:"id"`
	} `json:"installation"`
}

func GithubWebhook(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)

	if !utils.VerifyGithubSignature(r, body) {
		http.Error(w, "invalid signature", http.StatusUnauthorized)
		return
	}

	var payload GithubPushPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	if payload.Ref != "refs/heads/main" {
		fmt.Println("Ignoring non-main push")
		w.WriteHeader(200)
		return
	}

	appIDStr := os.Getenv("GITHUB_APP_ID")
	pemRaw := strings.ReplaceAll(os.Getenv("GITHUB_PRIVATE_KEY"), `\n`, "\n")

	appID, _ := strconv.ParseInt(appIDStr, 10, 64)

	installationID := payload.Installation.ID
	fmt.Println("installation id:", installationID)

	jwtToken, _ := utils.GenerateGithubJWT(appID, []byte(pemRaw))
	installationToken, err := utils.GetInstallationToken(jwtToken, installationID)
	if err != nil {
		fmt.Println("Error getting installation token:", err)
		http.Error(w, "failed to get installation token", 500)
		return
	}

	go utils.TriggerDeploy(payload.Repository.CloneURL, payload.After, "MODEL-ID-PLACEHOLDER", installationToken)

	w.WriteHeader(200)
}
