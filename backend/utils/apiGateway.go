package utils

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/apigatewayv2"
	"github.com/aws/aws-sdk-go-v2/service/apigatewayv2/types"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

func CreateAPIGatewayForLambda(ctx context.Context, lambdaArn string) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return "", err
	}

	client := apigatewayv2.NewFromConfig(cfg)

	// 1. Create API
	api, err := client.CreateApi(ctx, &apigatewayv2.CreateApiInput{
		Name:         aws.String("model-api"),
		ProtocolType: types.ProtocolTypeHttp,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create API: %w", err)
	}

	// 2. Create Integration
	integration, err := client.CreateIntegration(ctx, &apigatewayv2.CreateIntegrationInput{
		ApiId:                api.ApiId,
		IntegrationType:      types.IntegrationTypeAwsProxy,
		IntegrationUri:       aws.String(lambdaArn),
		PayloadFormatVersion: aws.String("2.0"),
	})
	if err != nil {
		return "", fmt.Errorf("failed to create integration: %w", err)
	}

	// 3. Create Route
	_, err = client.CreateRoute(ctx, &apigatewayv2.CreateRouteInput{
		ApiId:    api.ApiId,
		RouteKey: aws.String("POST /invoke"),
		Target:   aws.String("integrations/" + *integration.IntegrationId),
	})
	if err != nil {
		return "", fmt.Errorf("failed to create route: %w", err)
	}

	// 4. Allow API Gateway to invoke Lambda
	lc := lambda.NewFromConfig(cfg)

	accountID := "051826726578" // Replace with your account ID
	stage := "$default"
	sourceArn := fmt.Sprintf(
		"arn:aws:execute-api:%s:%s:%s/%s/POST/invoke",
		cfg.Region,
		accountID,
		*api.ApiId,
		stage,
	)

	_, err = lc.AddPermission(ctx, &lambda.AddPermissionInput{
		Action:       aws.String("lambda:InvokeFunction"),
		FunctionName: aws.String(lambdaArn),
		Principal:    aws.String("apigateway.amazonaws.com"),
		StatementId:  aws.String("apigw-invoke-" + *api.ApiId),
		SourceArn:    aws.String(sourceArn),
	})
	if err != nil {
		return "", fmt.Errorf("failed to add lambda permission: %w", err)
	}

	// 5. Deploy API
	_, err = client.CreateDeployment(ctx, &apigatewayv2.CreateDeploymentInput{
		ApiId: api.ApiId,
	})
	if err != nil {
		return "", fmt.Errorf("failed to deploy API: %w", err)
	}

	invokeURL := fmt.Sprintf(
		"https://%s.execute-api.%s.amazonaws.com/invoke",
		*api.ApiId,
		cfg.Region,
	)

	return invokeURL, nil
}
