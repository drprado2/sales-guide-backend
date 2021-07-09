package sns

import (
	"context"
	"testing"
)

func TestPublish(t *testing.T) {
	ctx := context.Background()
	if err := Setup(ctx); err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}
	msg := "This message coming from unit test"
	topicArn := "arn:aws:sns:sa-east-1:000000000000:sns-test"
	msgId, err := PublishSvc(ctx, topicArn, msg)
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}
	if msgId == "" {
		t.Errorf("expected message id not be empty")
	}
}
