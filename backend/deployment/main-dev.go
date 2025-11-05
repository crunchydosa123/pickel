package main_dev

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"

	"pickel-backend/routes"
)

var adapter *httpadapter.HandlerAdapter

func init() {
	router := routes.SetupRoutes()
	adapter = httpadapter.New(router)
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return adapter.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(handler)
}
