package sm

import (
	"context"
	"testing"
)

func TestGetValueFromKeySvc(t *testing.T) {
	ctx := context.Background()
	if err := Setup(ctx); err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}
	val, err := GetValueFromKeySvc(ctx, "test-sm", "secret-key")
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}
	if val != "test-key" {
		t.Errorf("expected error to be test-key, got %s", val)
	}
}
