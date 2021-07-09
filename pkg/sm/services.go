package sm

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/drprado2/react-redux-typescript/pkg/awsconfig"
)

var (
	client *secretsmanager.Client
)

func Setup(ctx context.Context) error {
	cfg, err := awsconfig.GetDefault(ctx)
	if err != nil {
		return err
	}

	client = secretsmanager.NewFromConfig(cfg)
	return nil
}

func GetValueFromKeySvc(ctx context.Context, secret string, key string) (string, error) {
	ipt := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secret),
		VersionStage: aws.String("AWSCURRENT"),
	}
	sec, err := client.GetSecretValue(ctx, ipt)
	if err != nil {
		return "", err
	}
	var secretString, decodedBinarySecret string
	var keysRes map[string]string

	if sec.SecretString != nil {
		secretString = *sec.SecretString
		json.Unmarshal([]byte(secretString), &keysRes)
	} else {
		decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(sec.SecretBinary)))
		len, err := base64.StdEncoding.Decode(decodedBinarySecretBytes, sec.SecretBinary)
		if err != nil {
			fmt.Println("Base64 Decode Error:", err)
		}
		decodedBinarySecret = string(decodedBinarySecretBytes[:len])
		json.Unmarshal([]byte(decodedBinarySecret), &keysRes)
	}

	return keysRes[key], nil
}
