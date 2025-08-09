package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"go-api-template/internal/adapter/httpserver"
	"go-api-template/internal/adapter/repository/dynamo"
	"go-api-template/internal/port"
	"go-api-template/internal/usecase/example"
	"go-api-template/internal/usecase/health"
	"go-api-template/pkg/dynamodb"
	"go-api-template/pkg/logger/zaplogger"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	localstackEndpoint := os.Getenv("LOCALSTACK_ENDPOINT")
	if localstackEndpoint == "" {
		localstackEndpoint = "http://localhost:4566"
	}
	region := getenv("AWS_REGION", "us-east-1")

	// Logger
	zl, err := zaplogger.FromEnv()
	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	defer zl.Sync()

	// DynamoDB client (opcional)
	var ddbPort port.DynamoDBPort
	if os.Getenv("ENABLE_DYNAMODB") != "false" { // habilitado por padrão
		ddbClient, err := dynamodb.New(ctx, dynamodb.Options{Region: region, Endpoint: localstackEndpoint})
		if err != nil {
			zl.Warn("failed to init dynamodb client", "error", err)
		} else {
			ddbPort = ddbClient
		}
	}

	healthSvc := health.NewServiceWithDeps(localstackEndpoint, ddbPort, zl)

	// Example CRUD service (opcional)
	var router http.Handler
	if ddbPort != nil {
		repo := dynamo.NewExampleRepository(ddbPort, getenv("EXAMPLE_TABLE", "example-items"))
		exampleSvc := example.NewService(repo, zl)
		router = httpserver.NewRouter(healthSvc, exampleSvc)
		zl.Info("example CRUD routes enabled", "table", repo.TableName())
	} else {
		router = httpserver.NewRouter(healthSvc)
	}

	addr := ":8080"
	zl.Info("starting api server", "addr", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		zl.Error("failed to start server", "error", err)
		log.Fatalf("failed to start server: %v", err)
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
