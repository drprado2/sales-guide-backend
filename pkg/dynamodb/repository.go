package dynamodb

import (
	"context"
	"errors"
)

type (
	Repository struct {
	}

	ListTables func(ctx context.Context) (string, error)
	GetItem    func(ctx context.Context, table string, keys map[string]interface{}, projectionExpression string, target interface{}) error
)

var (
	ItemNotFoundError = errors.New("item not found")
)
