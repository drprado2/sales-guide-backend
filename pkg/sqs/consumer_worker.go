package sqs

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/drprado2/react-redux-typescript/pkg/logs"
	"strconv"
	"time"
)

type (
	ConsumerWorker struct {
		client Client
	}

	AddConsumerInput struct {
		QueueName                    string
		ConcurrencyWorkers           int
		PoolingSeconds               int32
		MessagesPerPool              int32
		ResetVisibilityTimeoutOnFail bool
		Handler                      func(ctx context.Context, message types.Message) error
	}
)

func (cw *ConsumerWorker) StartConsumer(ctx context.Context, input *AddConsumerInput) error {
	logger := logs.Logger(ctx)
	qUrl, err := cw.client.GetQueueUrl(ctx, input.QueueName)
	if err != nil {
		return err
	}
	receiversChannels := make(chan types.Message, input.ConcurrencyWorkers)
	for i := 1; i <= input.ConcurrencyWorkers; i++ {
		go func() {
			for {
				select {
				case <-ctx.Done():
					logger.Warnf("Closing SQS consumer by ctx done, err=%v", ctx.Err())
					return
				case msg := <-receiversChannels:
					currentCtx := context.WithValue(ctx, "queue_name", input.QueueName)
					cidFromMsg := msg.MessageAttributes["cid"].StringValue
					visibilityTimeout := *msg.MessageAttributes["visibility_timeout"].StringValue
					vt, _ := strconv.Atoi(visibilityTimeout)
					currentCtx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(vt))
					logger := logs.Logger(currentCtx)
					if cidFromMsg != nil {
						currentCtx = context.WithValue(ctx, "cid", cidFromMsg)
					}
					err := input.Handler(currentCtx, msg)
					if err != nil {
						logger.Errorf("error consuming message queue=%s, err=%v", input.QueueName, err)
						if input.ResetVisibilityTimeoutOnFail {
							cw.client.ChangeMsgVisibilityTimeout(currentCtx, qUrl, *msg.ReceiptHandle, 0)
						}
					} else {
						cw.client.DeleteMsg(currentCtx, qUrl, *msg.ReceiptHandle)
					}
					cancel()
				}
			}
		}()
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				logger.Warnf("Closing SQS consumer by ctx done, err=%v", ctx.Err())
				return
			default:
				msgs, err := cw.client.GetMessages(ctx, qUrl, input.PoolingSeconds, input.MessagesPerPool)
				if err != nil {
					logs.Logger(ctx).Errorf("error gettings messages err=%v", err)
					continue
				}
				if len(msgs) == 0 {
					continue
				}
				for _, m := range msgs {
					receiversChannels <- m
				}
			}
		}
	}()
	return nil
}
