package kms

import "context"

type (
	Client struct {
		Decrypt
		Encrypt
	}

	Decrypt func(ctx context.Context, cipherTxt string) (string, error)
	Encrypt func(ctx context.Context, planTxt string, keyId string) (string, error)
)
