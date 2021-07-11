package s3

import (
	"bytes"
	"context"
	"testing"
	"time"
)

func TestFullFeatures(t *testing.T) {
	ctx := context.Background()
	if err := Setup(ctx); err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}

	bucketName := "sandpit-sample"
	testDir := "test-dir"
	testFileName := "myfile.txt"
	content := []byte("This is the content")

	key, err := PutFileSvc(ctx, testDir, testFileName, bucketName, content, nil)
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}

	files, err := GetFilesSvc(ctx, bucketName, 10, &testDir)
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}
	if len(files) != 1 {
		t.Errorf("expected to have 1 item, got %v", len(files))
	}
	if *files[0].Key != key {
		t.Errorf("expected to be %s item, got %s", key, *files[0].Key)
	}

	conRes, err := GetFileSvc(ctx, bucketName, key)
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}
	if bytes.Compare(conRes, content) != 0 {
		t.Errorf("expected content to be %s, got %s", string(content), string(conRes))
	}
	err = DeleteFileSvc(ctx, bucketName, key)
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}
	time.Sleep(time.Millisecond * 600)
	files, err = GetFilesSvc(ctx, bucketName, 10, &testDir)
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}
	if len(files) > 0 {
		t.Errorf("expected to have 0 item, got %v", len(files))
	}
}
