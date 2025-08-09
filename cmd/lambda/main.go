package main

import (
	"go-api-template/internal/handler"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler.LambdaHandler)
}
