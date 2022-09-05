package dynamodb

import (
	"context"
	"github.com/drprado2/sales-guide/pkg/instrumentation/logs"
	"github.com/google/uuid"
	"testing"
)

func TestListTablesSvc(t *testing.T) {
	ctx := context.Background()
	logs.Setup()
	if err := Setup(ctx); err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}
	tables, err := ListTablesSvc(ctx)
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}
	if len(tables) != 1 {
		t.Errorf("expected to have 1 table, got %v", len(tables))
	}
}

type TestModel struct {
	UserId    string
	GameTitle string
	TopScore  float32
	Age       int
}

func TestGetItemSvc(t *testing.T) {
	ctx := context.Background()
	logs.Setup()
	if err := Setup(ctx); err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}

	table := "test_table"
	keys := make(map[string]interface{})
	keys["UserId"] = "1"
	var result TestModel
	err := GetItemSvc(ctx, table, keys, "UserId, TopScore, GameTitle, Age", &result)
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}
	if result.UserId != "1" {
		t.Errorf("expected ID to be 1, got %v", result.UserId)
	}
	if result.TopScore != 23 {
		t.Errorf("expected TopScore to be 23, got %v", result.TopScore)
	}
	if result.GameTitle != "Test Game" {
		t.Errorf("expected GameTitle to be Test Game, got %v", result.GameTitle)
	}
	if result.Age != 27 {
		t.Errorf("expected Age to be 27, got %v", result.Age)
	}
}

func TestPutITemSvc_DeleteItemSvc(t *testing.T) {
	ctx := context.Background()
	logs.Setup()
	if err := Setup(ctx); err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}

	item := TestModel{
		UserId:    uuid.NewString(),
		GameTitle: "My Test Item",
		TopScore:  33,
		Age:       99,
	}

	table := "test_table"

	err := PutItemSvc(ctx, table, item)
	if err != nil {
		t.Fatalf("expected error to be nil, got %v", err)
	}

	keys := make(map[string]interface{})
	keys["UserId"] = item.UserId
	var result TestModel
	err = GetItemSvc(ctx, table, keys, "UserId, TopScore, GameTitle, Age", &result)
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}
	if result.UserId != item.UserId {
		t.Errorf("expected ID to be 1, got %v", result.UserId)
	}
	if result.TopScore != item.TopScore {
		t.Errorf("expected TopScore to be 23, got %v", result.TopScore)
	}
	if result.GameTitle != item.GameTitle {
		t.Errorf("expected GameTitle to be Test Game, got %v", result.GameTitle)
	}
	if result.Age != item.Age {
		t.Errorf("expected Age to be 27, got %v", result.Age)
	}

	updates := make(map[string]interface{})
	updates["Age"] = 166
	err = UpdateItemSvc(ctx, table, keys, updates)
	if err != nil {
		t.Fatalf("expected error to be nil, got %v", err)
	}
	err = GetItemSvc(ctx, table, keys, "UserId, TopScore, GameTitle, Age", &result)
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}
	if result.Age != 166 {
		t.Errorf("expected Age to be 166, got %v", result.Age)
	}

	err = DeleteItemSvc(ctx, table, keys)
	if err != nil {
		t.Fatalf("expected error to be nil, got %v", err)
	}

	err = GetItemSvc(ctx, table, keys, "UserId, TopScore, GameTitle, Age", &result)
	if err != ItemNotFoundError {
		t.Errorf("expected error to be %v, got %v", ItemNotFoundError, err)
	}
}

func TestBatchGetItemSvc(t *testing.T) {
	ctx := context.Background()
	logs.Setup()
	if err := Setup(ctx); err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}

	table := "test_table"
	filters := make([]map[string]interface{}, 2)

	keys1 := make(map[string]interface{})
	keys1["UserId"] = "1"
	keys2 := make(map[string]interface{})
	keys2["UserId"] = "2"
	filters[0] = keys1
	filters[1] = keys2
	var result TestModel
	results, err := BatchGetItemSvc(ctx, table, filters, "UserId, TopScore, GameTitle, Age", &result)
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected to have 2 items, got %v", len(results))
	}
	if results[0].(TestModel).UserId != "1" {
		t.Errorf("expected ID to be 1, got %v", results[0].(TestModel).UserId)
	}
	if results[0].(TestModel).TopScore != 23 {
		t.Errorf("expected TopScore to be 23, got %v", results[0].(TestModel).TopScore)
	}
	if results[0].(TestModel).GameTitle != "Test Game" {
		t.Errorf("expected GameTitle to be Test Game, got %v", results[0].(TestModel).GameTitle)
	}
	if results[0].(TestModel).Age != 27 {
		t.Errorf("expected Age to be 27, got %v", results[0].(TestModel).Age)
	}
	if results[1].(TestModel).UserId != "2" {
		t.Errorf("expected ID to be 2, got %v", results[1].(TestModel).UserId)
	}
	if results[1].(TestModel).TopScore != 43 {
		t.Errorf("expected TopScore to be 43, got %v", results[1].(TestModel).TopScore)
	}
	if results[1].(TestModel).GameTitle != "Test Game 2" {
		t.Errorf("expected GameTitle to be Test Game 2, got %v", results[1].(TestModel).GameTitle)
	}
	if results[1].(TestModel).Age != 33 {
		t.Errorf("expected Age to be 33, got %v", results[1].(TestModel).Age)
	}
}

func TestBatchPutItemsSvc(t *testing.T) {
	ctx := context.Background()
	logs.Setup()
	if err := Setup(ctx); err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}
	table := "test_table"

	items := make([]interface{}, 2)
	items[0] = &TestModel{
		UserId:    uuid.NewString(),
		GameTitle: "Item Test One",
		TopScore:  19,
		Age:       29,
	}
	items[1] = &TestModel{
		UserId:    uuid.NewString(),
		GameTitle: "Item Test Two",
		TopScore:  48,
		Age:       38,
	}

	err := BatchPutItemsSvc(ctx, table, items)

	filters := make([]map[string]interface{}, 2)

	keys1 := make(map[string]interface{})
	keys1["UserId"] = items[0].(*TestModel).UserId
	keys2 := make(map[string]interface{})
	keys2["UserId"] = items[1].(*TestModel).UserId
	filters[0] = keys1
	filters[1] = keys2
	var result TestModel
	results, err := BatchGetItemSvc(ctx, table, filters, "UserId, TopScore, GameTitle, Age", &result)
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected to have 2 items, got %v", len(results))
	}

	if err := BatchDeleteItemsSvc(ctx, table, filters); err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}
}
