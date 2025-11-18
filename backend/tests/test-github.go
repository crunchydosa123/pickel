package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func loadPrivateKey(path string) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed reading private key: %w", err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block")
	}

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSA private key: %w", err)
	}

	return key, nil
}

func generateJWT(appID int64, privateKey *rsa.PrivateKey) (string, error) {
	now := time.Now().Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iat": now,                      // issued at
		"exp": now + 540,                // expires in 9 minutes
		"iss": fmt.Sprintf("%d", appID), // GitHub App ID
	})

	signed, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed signing JWT: %w", err)
	}

	return signed, nil
}

func getInstallations(jwt string) error {
	req, err := http.NewRequest("GET", "https://api.github.com/app/installations", nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+jwt)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode >= 300 {
		return fmt.Errorf("GitHub API error (%d): %s", resp.StatusCode, string(body))
	}

	fmt.Println("Installations:")
	fmt.Println(string(body))
	return nil
}

func main() {
	// TODO: fill these
	const APP_ID = 2305067 // <-- your GitHub App ID
	privateKeyPath := "/Users/prathamgadkari/Projects/pickel/backend/pickel-deploy-bot.2025-11-16.private-key.pem"

	privKey, err := loadPrivateKey(privateKeyPath)
	if err != nil {
		panic(err)
	}

	jwt, err := generateJWT(APP_ID, privKey)
	if err != nil {
		panic(err)
	}

	fmt.Println("Generated JWT:\n", jwt, "\n")

	if err := getInstallations(jwt); err != nil {
		panic(err)
	}
}
