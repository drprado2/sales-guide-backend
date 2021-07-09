package dynamodb

import (
	"context"
	"errors"
)

type (
	Repository struct {
		ListTables       ListTables
		GetItem          GetItem
		PutItem          PutItem
		DeleteItem       DeleteItem
		UpdateItem       UpdateItem
		BatchGetItem     BatchGetItem
		BatchPutItems    BatchPutItems
		BatchDeleteItems BatchDeleteItems
	}

	ListTables       func(ctx context.Context) (string, error)
	GetItem          func(ctx context.Context, table string, keys map[string]interface{}, projectionExpression string, target interface{}) error
	PutItem          func(ctx context.Context, table string, item interface{}) error
	DeleteItem       func(ctx context.Context, table string, keys map[string]interface{}) error
	UpdateItem       func(ctx context.Context, table string, keys map[string]interface{}, updates map[string]interface{}) error
	BatchGetItem     func(ctx context.Context, table string, filters []map[string]interface{}, projectionExpression string, target interface{}) ([]interface{}, error)
	BatchPutItems    func(ctx context.Context, table string, elements []interface{}) error
	BatchDeleteItems func(ctx context.Context, table string, keys []map[string]interface{}) error
)

var (
	ItemNotFoundError = errors.New("item not found")
)
