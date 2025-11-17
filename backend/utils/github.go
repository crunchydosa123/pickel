package utils

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateGithubJWT(appID int64, privateKey []byte) (string, error) {
	now := time.Now()

	claims := jwt.MapClaims{
		"iat": now.Unix(),
		"exp": now.Add(time.Minute * 10).Unix(),
		"iss": appID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	return token.SignedString(privateKey)
}

func GetInstallationToken(jwtToken string, installationID int64) (string, error) {
	url := fmt.Sprintf("https://api.github.com/app/installations/%d/access_tokens", installationID)

	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Set("Authorization", "Bearer "+jwtToken)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)
	fmt.Println("GitHub installation token response:", string(respBytes))

	if resp.StatusCode != 201 {
		return "", fmt.Errorf("failed to get installation token: %s", resp.Status)
	}

	var body struct {
		Token string `json:"token"`
	}
	json.Unmarshal(respBytes, &body)

	return body.Token, nil
}

func VerifyGithubSignature(r *http.Request, body []byte) bool {
	sig := r.Header.Get("X-Hub-Signature-256")
	if sig == "" {
		return false
	}

	var githubWebhookSecret = os.Getenv("GITHUB_WEBHOOK_SECRET")

	mac := hmac.New(sha256.New, []byte(githubWebhookSecret))
	mac.Write(body)
	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(sig), []byte(expected))
}

func TriggerDeploy(repoURL, commitSHA, modelID string, installationToken string) error {
	tmpDir, err := os.MkdirTemp("", "repo-*")
	if err != nil {
		return fmt.Errorf("failed to created tmp dir %w", err)
	}
	defer os.RemoveAll(tmpDir)

	authedURL := strings.Replace(repoURL, "https://", "https://x-access-token:"+installationToken+"@", 1)

	fmt.Println("Cloning repo: ", authedURL)
	cmd := exec.Command("git", "clone", "--depth", "1", "--branch", "main", authedURL, tmpDir)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git clone failed: %v, %s", err, string(output))
	}

	cmd = exec.Command("git", "checkout", commitSHA)
	cmd.Dir = tmpDir
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git checkout failed %v, %s", err, string(output))
	}

	modelFilePath := filepath.Join(tmpDir, "model.py")
	if _, err := os.Stat(modelFilePath); err != nil {
		return fmt.Errorf("model file not found %w", err)
	}

	file, err := os.Open(modelFilePath)
	if err != nil {
		return fmt.Errorf("failed to open model file %w", err)
	}
	defer file.Close()

	s3Key := fmt.Sprintf("%s/%s", modelID, filepath.Base(modelFilePath))
	s3URL, err := UploadFileToS3(context.Background(), file, s3Key)
	if err != nil {
		return fmt.Errorf("S3 upload failed %w", err)
	}

	fmt.Println("S3 uploaded to: ", s3URL)

	lambdaArn, err := DeployToLambda(filepath.Base(modelFilePath), s3Key)
	if err != nil {
		return fmt.Errorf("API Gateway creation failed: %w", err)
	}

	apiURL, err := CreateAPIGatewayForLambda(context.Background(), lambdaArn)
	if err != nil {
		return fmt.Errorf("API Gateway creation failed: %w", err)
	}

	fmt.Println("Model deployed at API URL: ", apiURL)

	return nil

}
