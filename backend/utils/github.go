package utils

import (
	"context"
	"crypto/hmac"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
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

func GenerateGithubJWT(appID int64, pemBytes []byte) (string, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return "", fmt.Errorf("failed to decode PEM block")
	}

	var privateKey *rsa.PrivateKey
	var err error

	// Try PKCS1 first
	privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		// Try PKCS8
		key, err2 := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err2 != nil {
			return "", fmt.Errorf("failed parsing key: PKCS1=%v, PKCS8=%v", err, err2)
		}
		privateKey = key.(*rsa.PrivateKey)
	}

	now := time.Now().Unix()
	claims := jwt.MapClaims{
		"iat": now,
		"exp": now + 540,
		"iss": fmt.Sprintf("%d", appID),
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
	if installationToken == "" {
		return fmt.Errorf("empty installation token")
	}

	tmpDir, err := os.MkdirTemp("", "repo-*")
	if err != nil {
		return fmt.Errorf("failed to create tmp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	authedURL := strings.Replace(repoURL,
		"https://",
		"https://x-access-token:"+installationToken+"@",
		1,
	)

	fmt.Println("Cloning repo:", authedURL)

	// Clone into tmpDir
	cloneCmd := exec.Command("git", "clone", "--no-single-branch", "--depth=1", authedURL, tmpDir)
	cloneCmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	if out, err := cloneCmd.CombinedOutput(); err != nil {
		// try a full clone if shallow clone failed
		fmt.Printf("shallow clone failed: %v\n%s\nTrying full clone...\n", err, out)
		cloneCmd = exec.Command("git", "clone", "--no-single-branch", authedURL, tmpDir)
		cloneCmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
		if out2, err2 := cloneCmd.CombinedOutput(); err2 != nil {
			return fmt.Errorf("git clone failed (full): %v\n%s", err2, out2)
		}
	}

	// Try a few fetch strategies
	fmt.Println("Fetching commit:", commitSHA)

	tryFetch := func(args ...string) (string, error) {
		cmd := exec.Command("git", args...)
		cmd.Dir = tmpDir
		cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
		out, err := cmd.CombinedOutput()
		return string(out), err
	}

	// 1) direct fetch the SHA
	if out, err := tryFetch("fetch", "origin", commitSHA); err != nil {
		fmt.Printf("fetch origin %s failed: %v\n%s\n", commitSHA, err, out)

		// 2) fetch refs and all heads (slower)
		if out2, err2 := tryFetch("fetch", "--unshallow"); err2 == nil {
			fmt.Printf("unshallow output: %s\n", out2)
		} else {
			fmt.Printf("unshallow failed: %v\n%s\n", err2, out2)
		}

		// 3) try fetching all refs (fallback)
		if out3, err3 := tryFetch("fetch", "--no-tags", "origin"); err3 != nil {
			fmt.Printf("fetch origin (all) failed: %v\n%s\n", err3, out3)
		} else {
			fmt.Printf("fetch origin (all) succeeded: %s\n", out3)
		}
	} else {
		fmt.Println("fetch origin <sha> success")
	}

	// Checkout with a longer timeout
	fmt.Println("Checking out commit:", commitSHA)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	checkoutCmd := exec.CommandContext(ctx, "git", "checkout", commitSHA)
	checkoutCmd.Dir = tmpDir
	checkoutCmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")

	out, err := checkoutCmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		// timed out — try zip fallback
		fmt.Println("git checkout timed out, attempting zip fallback")
	} else if err != nil {
		fmt.Printf("git checkout failed: %v\n%s\n", err, out)
	} else {
		fmt.Println("git checkout success")
	}

	// If checkout failed or timed out, fallback to downloading zipball of the commit
	needZipFallback := (err != nil) || (ctx.Err() == context.DeadlineExceeded)
	if needZipFallback {
		fmt.Println("Falling back to download zipball for commit:", commitSHA)
		ownerRepo := strings.TrimPrefix(repoURL, "https://github.com/")
		// repoURL might include .git; strip it
		ownerRepo = strings.TrimSuffix(ownerRepo, ".git")
		parts := strings.Split(ownerRepo, "/")
		if len(parts) < 2 {
			return fmt.Errorf("unexpected repoURL format: %s", repoURL)
		}
		owner := parts[0]
		repo := parts[1]

		zipURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/zipball/%s", owner, repo, commitSHA)
		req, _ := http.NewRequest("GET", zipURL, nil)
		req.Header.Set("Authorization", "token "+installationToken)
		req.Header.Set("Accept", "application/vnd.github+json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return fmt.Errorf("failed to download zip: %w", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("zip download failed: %s — %s", resp.Status, string(body))
		}

		zipPath := filepath.Join(tmpDir, "repo.zip")
		f, err := os.Create(zipPath)
		if err != nil {
			return fmt.Errorf("failed to create zip file: %w", err)
		}
		if _, err := io.Copy(f, resp.Body); err != nil {
			f.Close()
			return fmt.Errorf("failed to write zip file: %w", err)
		}
		f.Close()

		// unzip into tmpDir (simple approach)
		unzipCmd := exec.Command("unzip", "-q", zipPath, "-d", tmpDir)
		if out, err := unzipCmd.CombinedOutput(); err != nil {
			return fmt.Errorf("unzip failed: %v\n%s", err, out)
		}

		// After unzip, find the single folder created (github zip creates a top-level dir)
		entries, err := os.ReadDir(tmpDir)
		if err != nil {
			return fmt.Errorf("reading tmp dir failed: %w", err)
		}
		var repoRoot string
		for _, e := range entries {
			if e.IsDir() {
				// skip if it's . or zip file; pick first dir that's not a temp file
				if strings.HasPrefix(e.Name(), "repo-") {
					continue
				}
				repoRoot = filepath.Join(tmpDir, e.Name())
				break
			}
		}
		if repoRoot == "" {
			// fallback: maybe unzip produced files in tmpDir directly
			repoRoot = tmpDir
		}
		// set tmpDir to repoRoot for remaining steps
		tmpDir = repoRoot
	}

	// check model file
	modelFilePath := filepath.Join(tmpDir, "model.py")
	if _, err := os.Stat(modelFilePath); err != nil {
		return fmt.Errorf("model file not found: %w", err)
	}

	file, err := os.Open(modelFilePath)
	if err != nil {
		return fmt.Errorf("failed to open model file: %w", err)
	}
	defer file.Close()

	s3Key := fmt.Sprintf("%s/%s", modelID, filepath.Base(modelFilePath))
	s3URL, err := UploadFileToS3(context.Background(), file, s3Key)
	if err != nil {
		return fmt.Errorf("S3 upload failed: %w", err)
	}
	fmt.Println("S3 uploaded to:", s3URL)

	lambdaArn, err := DeployToLambda(filepath.Base(modelFilePath), s3Key)
	if err != nil {
		return fmt.Errorf("Lambda deploy failed: %w", err)
	}

	apiURL, err := CreateAPIGatewayForLambda(context.Background(), lambdaArn)
	if err != nil {
		return fmt.Errorf("API Gateway creation failed: %w", err)
	}

	fmt.Println("Model deployed at:", apiURL)
	return nil
}
