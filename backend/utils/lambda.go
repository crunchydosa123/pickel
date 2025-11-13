package utils

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/google/uuid"
)

func DeployToLambda(fileName, s3Key string) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())

	if err != nil {
		return "", err
	}

	client := lambda.NewFromConfig(cfg)
	functionName := "model_" + uuid.New().String()

	input := &lambda.CreateFunctionInput{
		Code: &types.FunctionCode{
			S3Bucket: aws.String(os.Getenv("S3_BUCKET_NAME")),
			S3Key:    aws.String(s3Key),
		},
		FunctionName: aws.String(functionName),
		Handler:      aws.String("handler"),
		Role:         aws.String(os.Getenv("LAMBDA_ROLE_ARN")),
		Runtime:      types.RuntimePython39,
		Timeout:      aws.Int32(30),
		MemorySize:   aws.Int32(1024),
	}

	_, err = client.CreateFunction(context.Background(), input)
	if err != nil {
		return "", err
	}

	endpoint := fmt.Sprintf("https://%s.lambda-url.amazonaws.com/default/%s", os.Getenv("AWS_REGION"), functionName)
	return endpoint, nil

}
