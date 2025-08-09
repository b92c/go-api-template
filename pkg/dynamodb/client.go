package dynamodb

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	ddb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Client implementa a porta DynamoDBPort usando AWS SDK v2.
type Client struct {
	c *ddb.Client
}

type Options struct {
	Region   string
	Endpoint string // para LocalStack
}

func New(ctx context.Context, opt Options) (*Client, error) {
	var lopts []func(*config.LoadOptions) error
	if opt.Region != "" {
		lopts = append(lopts, config.WithRegion(opt.Region))
	}
	if opt.Endpoint != "" {
		lopts = append(lopts, config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: opt.Endpoint, HostnameImmutable: true}, nil
			},
		)))
	}
	cfg, err := config.LoadDefaultConfig(ctx, lopts...)
	if err != nil {
		return nil, err
	}
	return &Client{c: ddb.NewFromConfig(cfg)}, nil
}

func (cl *Client) Health(ctx context.Context) error {
	_, err := cl.c.ListTables(ctx, &ddb.ListTablesInput{Limit: aws.Int32(1)})
	return err
}

// PutItem com mapeamento simples (S/N/B/BOOL/M/L) automático a partir de tipos Go básicos.
func (cl *Client) PutItem(ctx context.Context, table string, item map[string]any) error {
	av, err := toAttributeValueMap(item)
	if err != nil {
		return err
	}
	_, err = cl.c.PutItem(ctx, &ddb.PutItemInput{TableName: aws.String(table), Item: av})
	return err
}

func (cl *Client) GetItem(ctx context.Context, table string, key map[string]any) (map[string]any, error) {
	av, err := toAttributeValueMap(key)
	if err != nil {
		return nil, err
	}
	out, err := cl.c.GetItem(ctx, &ddb.GetItemInput{TableName: aws.String(table), Key: av})
	if err != nil {
		return nil, err
	}
	if out.Item == nil {
		return nil, errors.New("not found")
	}
	return fromAttributeValueMap(out.Item)
}

func (cl *Client) DeleteItem(ctx context.Context, table string, key map[string]any) error {
	av, err := toAttributeValueMap(key)
	if err != nil {
		return err
	}
	_, err = cl.c.DeleteItem(ctx, &ddb.DeleteItemInput{TableName: aws.String(table), Key: av})
	return err
}

func (cl *Client) Scan(ctx context.Context, table string, limit int32) ([]map[string]any, error) {
	out, err := cl.c.Scan(ctx, &ddb.ScanInput{TableName: aws.String(table), Limit: aws.Int32(limit)})
	if err != nil {
		return nil, err
	}
	items := make([]map[string]any, 0, len(out.Items))
	for _, it := range out.Items {
		m, err := fromAttributeValueMap(it)
		if err != nil {
			return nil, err
		}
		items = append(items, m)
	}
	return items, nil
}

func toAttributeValueMap(m map[string]any) (map[string]types.AttributeValue, error) {
	res := make(map[string]types.AttributeValue, len(m))
	for k, v := range m {
		sw, err := toAttributeValue(v)
		if err != nil {
			return nil, err
		}
		res[k] = sw
	}
	return res, nil
}

func toAttributeValue(v any) (types.AttributeValue, error) {
	switch t := v.(type) {
	case string:
		return &types.AttributeValueMemberS{Value: t}, nil
	case []byte:
		return &types.AttributeValueMemberB{Value: t}, nil
	case int:
		return &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", t)}, nil
	case int32:
		return &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", t)}, nil
	case int64:
		return &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", t)}, nil
	case float32:
		return &types.AttributeValueMemberN{Value: fmt.Sprintf("%g", t)}, nil
	case float64:
		return &types.AttributeValueMemberN{Value: fmt.Sprintf("%g", t)}, nil
	case bool:
		return &types.AttributeValueMemberBOOL{Value: t}, nil
	case map[string]any:
		m, err := toAttributeValueMap(t)
		if err != nil {
			return nil, err
		}
		return &types.AttributeValueMemberM{Value: m}, nil
	case []any:
		arr := make([]types.AttributeValue, 0, len(t))
		for _, it := range t {
			av, err := toAttributeValue(it)
			if err != nil {
				return nil, err
			}
			arr = append(arr, av)
		}
		return &types.AttributeValueMemberL{Value: arr}, nil
	default:
		return nil, errors.New("unsupported type")
	}
}

func fromAttributeValueMap(m map[string]types.AttributeValue) (map[string]any, error) {
	res := map[string]any{}
	for k, v := range m {
		res[k] = fromAttributeValue(v)
	}
	return res, nil
}

func fromAttributeValue(v types.AttributeValue) any {
	switch t := v.(type) {
	case *types.AttributeValueMemberS:
		return t.Value
	case *types.AttributeValueMemberB:
		return t.Value
	case *types.AttributeValueMemberN:
		return t.Value
	case *types.AttributeValueMemberBOOL:
		return t.Value
	case *types.AttributeValueMemberM:
		m := map[string]any{}
		for k, v := range t.Value {
			m[k] = fromAttributeValue(v)
		}
		return m
	case *types.AttributeValueMemberL:
		arr := make([]any, 0, len(t.Value))
		for _, v := range t.Value {
			arr = append(arr, fromAttributeValue(v))
		}
		return arr
	default:
		return nil
	}
}
