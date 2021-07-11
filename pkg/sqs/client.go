package sqs

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"time"
)

type (
	Client struct {
		GetQueueUrl
		GetMessages
		ChangeMsgVisibilityTimeout
		DeleteMsg
		SendMsg
		SendJsonMsg
		SendFifoMsg
		SendJsonFifoMsg
		CreateVirtualQueue
		DeleteQueue
		PurgeQueue
	}

	GetQueueUrl                func(ctx context.Context, queueName string) (string, error)
	GetMessages                func(ctx context.Context, queueUrl string, waitTimeSeconds int32, maxNumberOfMessages int32) ([]types.Message, error)
	ChangeMsgVisibilityTimeout func(ctx context.Context, queueUrl string, receiptId string, value time.Duration) error
	DeleteMsg                  func(ctx context.Context, queueUrl string, receiptId string) error
	SendMsg                    func(ctx context.Context, queueUrl string, delaySeconds int32, msg string, msgAttributes map[string]interface{}) (string, error)
	SendJsonMsg                func(ctx context.Context, queueUrl string, delaySeconds int32, msg interface{}, msgAttributes map[string]interface{}) (string, error)
	SendFifoMsg                func(ctx context.Context, queueUrl string, msg string, msgGroupId string, msgAttributes map[string]interface{}) (string, error)
	SendJsonFifoMsg            func(ctx context.Context, queueUrl string, msg interface{}, msgGroupId string, msgAttributes map[string]interface{}) (string, error)
	CreateVirtualQueue         func(ctx context.Context, baseQueueName string, baseQueueUrl string, tempQueueName string) (string, error)
	DeleteQueue                func(ctx context.Context, queueUrl string) error
	PurgeQueue                 func(ctx context.Context, queueUrl string) error
)
