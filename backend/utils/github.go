package utils

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func TriggerDeploy(repoURL, commitSHA, modelID string) error {
	tmpDir, err := os.MkdirTemp("", "repo-*")
	if err != nil {
		return fmt.Errorf("failed to created tmp dir %w", err)
	}
	defer os.RemoveAll(tmpDir)

	fmt.Println("Cloning repo: ", repoURL)
	cmd := exec.Command("git", "clone", "--depth", "1", "--branch", "main", repoURL, tmpDir)
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
