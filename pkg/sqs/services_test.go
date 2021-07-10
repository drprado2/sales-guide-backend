package sqs

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"strconv"
	"testing"
	"time"
)

type (
	TestModel struct {
		Id   string `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
		Age  int    `json:"age,omitempty"`
	}
)

func TestFullFeatures(t *testing.T) {
	ctx := context.Background()
	if err := Setup(ctx); err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}
	normalQueueName := "test_queue"
	normalQueueUrl, err := GetQueueUrlSvc(ctx, normalQueueName)
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}

	fifoQueueName := "test_queue_fifo.fifo"
	fifoQueueUrl, err := GetQueueUrlSvc(ctx, fifoQueueName)
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}
	groupFifoId := uuid.NewString()

	PurgeQueueSvc(ctx, normalQueueUrl)
	PurgeQueueSvc(ctx, fifoQueueUrl)
	time.Sleep(time.Millisecond * 500)

	if msgs, err := GetMessagesSvc(ctx, normalQueueUrl, 0, 1); err != nil || len(msgs) > 0 {
		t.Errorf("Expected not found messages, msgs=%v, err=%v", msgs, err)
	}
	if msgs, err := GetMessagesSvc(ctx, fifoQueueUrl, 0, 1); err != nil || len(msgs) > 0 {
		t.Errorf("Expected not found messages, msgs=%v, err=%v", msgs, err)
	}

	textMsg := "Text msg"
	attr := make(map[string]interface{})
	attr["number_field"] = 33
	attr["string_field"] = "text"
	model := &TestModel{
		Id:   uuid.NewString(),
		Name: "Test Model",
		Age:  23,
	}
	attr["binary_field"] = model
	nMsgId1, err := SendMsgSvc(ctx, normalQueueUrl, 0, textMsg, attr)
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}
	nMsgId2, err := SendJsonMsgSvc(ctx, normalQueueUrl, 0, model, nil)
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}

	fMsgId1, err := SendFifoMsgSvc(ctx, fifoQueueUrl, textMsg, groupFifoId, attr)
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}
	fMsgId2, err := SendJsonFifoMsgSvc(ctx, fifoQueueUrl, model, groupFifoId, nil)
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}

	nMsgs, err := GetMessagesSvc(ctx, normalQueueUrl, 20, 2)
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}
	if len(nMsgs) != 2 {
		t.Errorf("expected to have 2 messages, got %v", len(nMsgs))
	}
	for _, m := range nMsgs {
		if *m.MessageId == nMsgId1 {
			if *m.Body != textMsg {
				t.Errorf("message text to be %s got %s", textMsg, *m.Body)
			}
			nf := *m.MessageAttributes["number_field"].StringValue
			numberField, err := strconv.Atoi(nf)
			if err != nil {
				t.Errorf("expected error to be nil, got %v", err)
			}
			if numberField != 33 {
				t.Errorf("expected number field to be 33, got %v", numberField)
			}
			sf := *m.MessageAttributes["string_field"].StringValue
			if sf != "text" {
				t.Errorf("expected string field to be text, got %s", sf)
			}
			bf := m.MessageAttributes["binary_field"].BinaryValue
			var tempModel TestModel
			if err := json.Unmarshal(bf, &tempModel); err != nil {
				t.Errorf("expected error to be nil, got %v", err)
			}
			if tempModel.Id != model.Id || tempModel.Name != model.Name || tempModel.Age != model.Age {
				t.Errorf("expected message to be equal msg=%v, model=%v", tempModel, model)
			}
		} else {
			temp := &TestModel{}
			if *m.MessageId != nMsgId2 {
				t.Errorf("expected message ID to be %s, got %s", nMsgId2, *m.MessageId)
			}
			if err := json.Unmarshal([]byte(*m.Body), temp); err != nil {
				t.Errorf("expected error to be nil, got %v", err)
			}
			if temp.Id != model.Id || temp.Name != model.Name || temp.Age != model.Age {
				t.Errorf("expected message to be equal msg=%v, model=%v", temp, model)
			}
		}
	}

	fMsgs, err := GetMessagesSvc(ctx, fifoQueueUrl, 0, 2)
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}
	if len(fMsgs) != 2 {
		t.Errorf("expected to have 2 messages, got %v", len(fMsgs))
	}
	if *fMsgs[0].MessageId != fMsgId1 {
		t.Errorf("expected first fifo message to be %s, got %s", fMsgId1, *fMsgs[0].MessageId)
	}
	if *fMsgs[1].MessageId != fMsgId2 {
		t.Errorf("expected second fifo message to be %s, got %s", fMsgId1, *fMsgs[0].MessageId)
	}

	ChangeMsgVisibilityTimeoutSvc(ctx, normalQueueUrl, *nMsgs[0].ReceiptHandle, 0)
	tempMsgs, _ := GetMessagesSvc(ctx, normalQueueUrl, 20, 1)
	if len(tempMsgs) != 1 {
		t.Errorf("expected to have 1 message got %v", len(tempMsgs))
	}

	if err := DeleteMsgSvc(ctx, normalQueueUrl, *tempMsgs[0].ReceiptHandle); err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}
	if err := DeleteMsgSvc(ctx, normalQueueUrl, *nMsgs[1].ReceiptHandle); err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}
	if err := DeleteMsgSvc(ctx, fifoQueueUrl, *fMsgs[0].ReceiptHandle); err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}
	if err := DeleteMsgSvc(ctx, fifoQueueUrl, *fMsgs[1].ReceiptHandle); err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}
}

func TestVirtualQueue(t *testing.T) {
	ctx := context.Background()
	if err := Setup(ctx); err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}
	normalQueueName := "test_queue"
	normalQueueUrl, err := GetQueueUrlSvc(ctx, normalQueueName)
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}
	PurgeQueueSvc(ctx, normalQueueUrl)
	tempQueueName := "testTemp_2"

	qName, err := CreateVirtualQueueSvc(ctx, normalQueueName, normalQueueUrl, tempQueueName)
	if err != nil {
		t.Fatalf("expected error to be nil, got %v", err)
	}
	defer DeleteQueueSvc(ctx, qName)

	go func() {
		for {
			msgs, err := GetMessagesSvc(ctx, normalQueueUrl, 20, 1)
			if err != nil {
				continue
			}
			if len(msgs) == 0 {
				continue
			}
			respQueue := *msgs[0].MessageAttributes["ResponseQueueUrl"].StringValue
			_, err = SendMsgSvc(ctx, respQueue, 0, "receive", nil)
			if err != nil {
				t.Fatalf("error publishing msg err=%v", err)
			}
		}
	}()

	SendMsgSvc(ctx, normalQueueUrl, 0, "come back", map[string]interface{}{
		"ResponseQueueUrl": qName,
	})
	time.Sleep(time.Millisecond * 600)
	finalMsg, err := GetMessagesSvc(ctx, qName, 20, 1)
	if err != nil {
		t.Fatalf("expected error to be nil, got %v", err)
	}
	if len(finalMsg) != 1 {
		t.Fatalf("expected to find 1 message")
	}
	if *finalMsg[0].Body != "receive" {
		t.Errorf("expected msg content to be receive, got %s", *finalMsg[0].Body)
	}
}
