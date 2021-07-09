package redis

import (
	"context"
	"github.com/google/uuid"
	"testing"
)

type TestObj struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

func TestSetKeySvc_GetKey_DeleteKey(t *testing.T) {
	ctx := context.Background()
	Setup(ctx)

	val1 := "test"
	val2 := TestObj{
		Id:   uuid.NewString(),
		Name: "Test name",
	}

	if err := SetKeySvc(ctx, "key1", val1, 0); err != nil {
		t.Fatalf("expected error to be nil, got %v", err)
	}
	if err := SetKeyJsonSvc(ctx, "key2", val2, 0); err != nil {
		t.Fatalf("expected error to be nil, got %v", err)
	}

	var res2 TestObj

	res1, err := GetFromKeySvc(ctx, "key1")
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}
	if res1 == nil {
		t.Errorf("expected to find key1")
	}
	ok, err := UnmarshalFromKeySvc(ctx, "key2", &res2)
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}
	if !ok {
		t.Errorf("expected to find key2")
	}

	if *res1 != val1 {
		t.Errorf("expected value1 to be %v, got %v", res1, val1)
	}
	if res2 != val2 {
		t.Errorf("expected value2 to be %v, got %v", res2, val2)
	}

	DeleteKeySvc(ctx, "key1")
	DeleteKeySvc(ctx, "key2")
	res1, err = GetFromKeySvc(ctx, "key1")
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}
	if res1 != nil {
		t.Errorf("expected key1 to be nil, got %v", res1)
	}
	ress2, err := GetFromKeySvc(ctx, "key2")
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}
	if ress2 != nil {
		t.Errorf("expected key2 to be nil, got %v", res1)
	}
}
