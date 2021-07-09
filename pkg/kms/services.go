package kms

import (
	"context"
	b64 "encoding/base64"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/drprado2/react-redux-typescript/pkg/awsconfig"
)

var (
	client *kms.Client
)

func Setup(ctx context.Context) error {
	cfg, err := awsconfig.GetDefault(ctx)
	if err != nil {
		return err
	}

	client = kms.NewFromConfig(cfg)
	return nil
}

func DecryptSvc(ctx context.Context, cipherTxt string) (string, error) {
	blob, err := b64.StdEncoding.DecodeString(cipherTxt)
	if err != nil {
		return "", err
	}

	ipt := &kms.DecryptInput{
		CiphertextBlob: blob,
	}

	result, err := client.Decrypt(ctx, ipt)
	if err != nil {
		return "", err
	}
	return string(result.Plaintext), nil
}

func EncryptSvc(ctx context.Context, planTxt string, keyId string) (string, error) {
	input := &kms.EncryptInput{
		KeyId:     aws.String(keyId),
		Plaintext: []byte(planTxt),
	}

	result, err := client.Encrypt(ctx, input)
	if err != nil {
		return "", err
	}
	return b64.StdEncoding.EncodeToString(result.CiphertextBlob), nil
}
