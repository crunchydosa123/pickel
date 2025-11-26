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
	defer r.Body.Close()

	if !utils.VerifyGithubSignature(r, body) {
		http.Error(w, "invalid signature", http.StatusUnauthorized)
		return
	}

	event := r.Header.Get("X-GitHub-Event")
	switch event {
	case "installation_repositories":
		utils.HandleInstallationEvent(body)
	case "push":
		utils.HandlePushEvent(body)
	default:
		fmt.Println("Ignoring event:", event)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func GithubCallback(w http.ResponseWriter, r *http.Request) {
	stateEncoded := r.URL.Query().Get("state")
	if stateEncoded == "" {
		http.Error(w, "missing state", http.StatusBadRequest)
		return
	}

	var state struct {
		ModelID string `json:"model_id"`
		RepoURL string `json:"repo_url"`
	}
	err := json.Unmarshal([]byte(stateEncoded), &state)
	if err != nil {
		http.Error(w, "invalid state", http.StatusBadRequest)
		return
	}

	fmt.Println("Decoded state:", state)

	// TODO: Save model_id + repo_url temporarily
	w.Write([]byte("Installation complete, you can close this window"))
}

type CallbackResponse struct {
	RepoURL        string `json:"repo_url"`
	InstallationId int64  `json:"installation_id"`
}

func GithubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	//ctx := context.Background()
	body, _ := io.ReadAll(r.Body)
	fmt.Println("Webhook payload:", string(body))

	/*installationID := r.URL.Query().Get("installation_id")
	modelID := r.URL.Query().Get("model_id")

	if installationID == "" {
		http.Error(w, "missing installation_id", 400)
		return
	}
	if modelID == "" {
		http.Error(w, "missing model_id", 400)
		return
	}

	instID, err := strconv.ParseInt(installationID, 10, 64)
	if err != nil {
		http.Error(w, "invalid installation_id", 400)
		return
	}

	appIDStr := os.Getenv("GITHUB_APP_ID")
	pemRaw := strings.ReplaceAll(os.Getenv("GITHUB_PRIVATE_KEY"), `\n`, "\n")

	appID, _ := strconv.ParseInt(appIDStr, 10, 64)

	jwtToken, err := utils.GenerateGithubJWT(appID, []byte(pemRaw))
	if err != nil {
		http.Error(w, "failed to generate jwt", 500)
		return
	}

	token, err := utils.GetInstallationToken(jwtToken, instID)
	if err != nil {
		fmt.Println("Token error:", err)
		http.Error(w, "failed to create installation token", 500)
		return
	}

	repos, err := utils.FetchInstallationRepos(token)
	if err != nil {
		fmt.Println("Repo error:", err)
		http.Error(w, "failed to fetch repos", 500)
		return
	}

	if len(repos) == 0 {
		http.Error(w, "no repositories found", 400)
		return
	}

	repoURL := repos[0].HtmlUrl
	fmt.Println("Selected repo:", repoURL)

	db := utils.GetDB()
	_, err = db.Exec(
		ctx,
		"UPDATE modelcode SET github_repo_url = $1 WHERE model_id = $2",
		repoURL,
		modelID,
	)
	if err != nil {
		fmt.Println("DB error:", err)
		http.Error(w, "failed to update modelcode", 500)
		return
	}

	resp := CallbackResponse{
		RepoURL:        repoURL,
		InstallationId: instID,
	}*/

	w.Header().Set("Content-Type", "application/json")
	//json.NewEncoder(w).Encode(resp)
}

func InstalledReposHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: replace this with your logic to get installation ID for the current user/model
	installationIDStr := "96640194"
	if installationIDStr == "" {
		http.Error(w, "installation_id is required", http.StatusBadRequest)
		return
	}
	installationID, _ := strconv.ParseInt(installationIDStr, 10, 64)

	appIDStr := os.Getenv("GITHUB_APP_ID")
	pemRaw := strings.ReplaceAll(os.Getenv("GITHUB_PRIVATE_KEY"), `\n`, "\n")
	appID, _ := strconv.ParseInt(appIDStr, 10, 64)

	jwtToken, err := utils.GenerateGithubJWT(appID, []byte(pemRaw))
	if err != nil {
		http.Error(w, "failed to generate JWT: "+err.Error(), http.StatusInternalServerError)
		return
	}

	installationToken, err := utils.GetInstallationToken(jwtToken, installationID)
	if err != nil {
		http.Error(w, "failed to get installation token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	repos, err := utils.FetchInstallationRepos(installationToken)
	if err != nil {
		http.Error(w, "failed to fetch repos: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"repositories": repos,
	})
}
