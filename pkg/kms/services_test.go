package kms

import (
	"context"
	"testing"
)

func TestEncryptSvc_DecryptSvc(t *testing.T) {
	ctx := context.Background()
	if err := Setup(ctx); err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}
	planTxt := "testText"
	keyId := "45a1f0a2-1ec7-474f-be64-615b040b675b"
	cipher, err := EncryptSvc(ctx, planTxt, keyId)
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}
	result, err := DecryptSvc(ctx, cipher)
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}
	if result != planTxt {
		t.Errorf("expected result to be %s, got %s", planTxt, result)
	}
}
