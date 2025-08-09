package port

import "context"

// DynamoDBPort define operações necessárias para os casos de uso, com mapas simples de valores (sem tipos do SDK).
// Implementações concretas ficam em pkg/dynamodb.
type DynamoDBPort interface {
	PutItem(ctx context.Context, table string, item map[string]any) error
	GetItem(ctx context.Context, table string, key map[string]any) (map[string]any, error)
	DeleteItem(ctx context.Context, table string, key map[string]any) error
	Scan(ctx context.Context, table string, limit int32) ([]map[string]any, error)
	Health(ctx context.Context) error
}
