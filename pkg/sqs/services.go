package sqs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/drprado2/react-redux-typescript/pkg/awsconfig"
	"strconv"
	"time"
)

var (
	client *sqs.Client
)

func Setup(ctx context.Context) error {
	cfg, err := awsconfig.GetDefault(ctx)
	if err != nil {
		return err
	}

	client = sqs.NewFromConfig(cfg)
	return nil
}

func GetQueueUrlSvc(ctx context.Context, queueName string) (string, error) {
	ipt := &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	}

	urlResult, err := client.GetQueueUrl(ctx, ipt)
	if err != nil {
		return "", err
	}

	return *urlResult.QueueUrl, nil
}

func GetMessagesSvc(ctx context.Context, queueUrl string, waitTimeSeconds int32, maxNumberOfMessages int32) ([]types.Message, error) {
	ipt := &sqs.ReceiveMessageInput{
		QueueUrl: aws.String(queueUrl),
		AttributeNames: []types.QueueAttributeName{
			types.QueueAttributeNameAll,
		},
		MaxNumberOfMessages: maxNumberOfMessages,
		MessageAttributeNames: []string{
			string(types.QueueAttributeNameAll),
		},
		WaitTimeSeconds: waitTimeSeconds,
	}

	msgResult, err := client.ReceiveMessage(ctx, ipt)
	if err != nil {
		return nil, err
	}

	return msgResult.Messages, nil
}

func ChangeMsgVisibilityTimeoutSvc(ctx context.Context, queueUrl string, receiptId string, value time.Duration) error {
	ipt := &sqs.ChangeMessageVisibilityInput{
		ReceiptHandle:     aws.String(receiptId),
		QueueUrl:          aws.String(queueUrl),
		VisibilityTimeout: int32(value),
	}

	_, err := client.ChangeMessageVisibility(ctx, ipt)
	return err
}

func DeleteMsgSvc(ctx context.Context, queueUrl string, receiptId string) error {
	ipt := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueUrl),
		ReceiptHandle: aws.String(receiptId),
	}

	_, err := client.DeleteMessage(ctx, ipt)
	return err
}

func CastInterfaceToMsgAttrValue(ipt map[string]interface{}) (map[string]types.MessageAttributeValue, error) {
	result := make(map[string]types.MessageAttributeValue)
	if ipt == nil {
		return result, nil
	}
	for k, v := range ipt {
		var parse types.MessageAttributeValue
		switch v.(type) {
		case int:
			parse.DataType = aws.String("Number")
			parse.StringValue = aws.String(strconv.Itoa(v.(int)))
		case float32:
			parse.DataType = aws.String("Number")
			parse.StringValue = aws.String(fmt.Sprintf("%f", v))
		case float64:
			parse.DataType = aws.String("Number")
			parse.StringValue = aws.String(fmt.Sprintf("%f", v))
		case string:
			parse.DataType = aws.String("String")
			parse.StringValue = aws.String(v.(string))
		default:
			j, err := json.Marshal(v)
			if err != nil {
				return nil, err
			}
			parse.DataType = aws.String("Binary")
			parse.BinaryValue = j
		}
		result[k] = parse
	}
	return result, nil
}

func SendMsgSvc(ctx context.Context, queueUrl string, delaySeconds int32, msg string, msgAttributes map[string]interface{}) (string, error) {
	attrs, err := CastInterfaceToMsgAttrValue(msgAttributes)
	if err != nil {
		return "", err
	}
	ipt := &sqs.SendMessageInput{
		MessageBody:       aws.String(msg),
		QueueUrl:          aws.String(queueUrl),
		DelaySeconds:      delaySeconds,
		MessageAttributes: attrs,
	}

	resp, err := client.SendMessage(ctx, ipt)
	if err != nil {
		return "", err
	}
	return *resp.MessageId, err
}

func SendJsonMsgSvc(ctx context.Context, queueUrl string, delaySeconds int32, msg interface{}, msgAttributes map[string]interface{}) (string, error) {
	attrs, err := CastInterfaceToMsgAttrValue(msgAttributes)
	if err != nil {
		return "", err
	}
	cast, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}
	ipt := &sqs.SendMessageInput{
		MessageBody:       aws.String(string(cast)),
		QueueUrl:          aws.String(queueUrl),
		DelaySeconds:      delaySeconds,
		MessageAttributes: attrs,
	}

	resp, err := client.SendMessage(ctx, ipt)
	if err != nil {
		return "", err
	}
	return *resp.MessageId, err
}

func SendFifoMsgSvc(ctx context.Context, queueUrl string, msg string, msgGroupId string, msgAttributes map[string]interface{}) (string, error) {
	attrs, err := CastInterfaceToMsgAttrValue(msgAttributes)
	if err != nil {
		return "", err
	}
	ipt := &sqs.SendMessageInput{
		MessageBody:       aws.String(msg),
		QueueUrl:          aws.String(queueUrl),
		MessageAttributes: attrs,
		MessageGroupId:    aws.String(msgGroupId),
	}

	resp, err := client.SendMessage(ctx, ipt)
	if err != nil {
		return "", err
	}
	return *resp.MessageId, err
}

func SendJsonFifoMsgSvc(ctx context.Context, queueUrl string, msg interface{}, msgGroupId string, msgAttributes map[string]interface{}) (string, error) {
	attrs, err := CastInterfaceToMsgAttrValue(msgAttributes)
	if err != nil {
		return "", err
	}
	cast, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}
	ipt := &sqs.SendMessageInput{
		MessageBody:       aws.String(string(cast)),
		QueueUrl:          aws.String(queueUrl),
		MessageAttributes: attrs,
		MessageGroupId:    aws.String(msgGroupId),
	}

	resp, err := client.SendMessage(ctx, ipt)
	if err != nil {
		return "", err
	}
	return *resp.MessageId, err
}

func CreateVirtualQueueSvc(ctx context.Context, baseQueueName string, baseQueueUrl string, tempQueueName string) (string, error) {
	queueName := fmt.Sprintf("%s#%s", baseQueueName, tempQueueName)
	ipt := &sqs.CreateQueueInput{
		QueueName: aws.String(queueName),
		Attributes: map[string]string{
			"HostQueueURL": baseQueueUrl,
		},
	}
	resp, err := client.CreateQueue(ctx, ipt)
	return *resp.QueueUrl, err
}

func DeleteQueueSvc(ctx context.Context, queueUrl string) error {
	ipt := &sqs.DeleteQueueInput{
		QueueUrl: aws.String(queueUrl),
	}
	_, err := client.DeleteQueue(ctx, ipt)
	return err
}

func PurgeQueueSvc(ctx context.Context, queueUrl string) error {
	ipt := &sqs.PurgeQueueInput{
		QueueUrl: aws.String(queueUrl),
	}
	_, err := client.PurgeQueue(ctx, ipt)
	return err
}
