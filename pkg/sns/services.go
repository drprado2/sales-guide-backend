package sns

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/drprado2/sales-guide/pkg/awsconfig"
)

var (
	client *sns.Client
)

func Setup(ctx context.Context) error {
	cfg, err := awsconfig.GetDefault(ctx)
	if err != nil {
		return err
	}

	client = sns.NewFromConfig(cfg)
	return nil
}

func PublishSvc(ctx context.Context, topicArn string, msg string) (string, error) {
	input := &sns.PublishInput{
		Message:  aws.String(msg),
		TopicArn: aws.String(topicArn),
	}

	result, err := client.Publish(ctx, input)
	if err != nil {
		return "", err
	}
	return *result.MessageId, nil
}
