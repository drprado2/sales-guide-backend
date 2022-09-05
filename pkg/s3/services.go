package s3

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/drprado2/sales-guide/configs"
	"github.com/drprado2/sales-guide/pkg/awsconfig"
	"io/ioutil"
	"time"
)

var (
	client *s3.Client
)

func Setup(ctx context.Context) error {
	envs := configs.Get()
	cfg, err := awsconfig.GetDefault(ctx)
	if err != nil {
		return err
	}

	client = s3.NewFromConfig(cfg, func(options *s3.Options) {
		options.UsePathStyle = envs.ForceS3PathStyle
	})
	return nil
}

func PutFileSvc(ctx context.Context, dirName string, fileName string, bucket string, content []byte, expires *time.Time) (string, error) {
	fileKey := fileName
	if dirName != "" {
		fileKey = fmt.Sprintf("%s/%s", dirName, fileName)
	}
	reader := bytes.NewReader(content)
	input := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileKey),
		Body:   reader,

		Expires: expires,
	}
	_, err := client.PutObject(ctx, input)
	return fileKey, err
}

func GetFilesSvc(ctx context.Context, bucket string, maxQuantity int32, prefix *string) ([]types.Object, error) {
	input := &s3.ListObjectsV2Input{
		Bucket:  aws.String(bucket),
		MaxKeys: maxQuantity,
		Prefix:  prefix,
	}
	res, err := client.ListObjectsV2(ctx, input)
	if err != nil {
		return nil, err
	}
	return res.Contents, nil
}

func GetFileSvc(ctx context.Context, bucket string, fileKey string) ([]byte, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileKey),
	}
	res, err := client.GetObject(ctx, input)
	if err != nil {
		return nil, err
	}
	obj, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func DeleteFileSvc(ctx context.Context, bucket string, fileKey string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileKey),
	}
	_, err := client.DeleteObject(ctx, input)
	return err
}
