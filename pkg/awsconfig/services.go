package awsconfig

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/drprado2/sales-guide/configs"
	"github.com/drprado2/sales-guide/pkg/instrumentation/logs"
)

func GetDefault(ctx context.Context) (aws.Config, error) {
	envs := configs.Get()
	customEndpointResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		if envs.AwsEndpoint != "" {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           envs.AwsEndpoint,
				SigningRegion: envs.AwsRegion,
			}, nil
		}

		// returning EndpointNotFoundError will allow the service to fallback to it's default resolution
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	// Get configs from env var
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithLogger(logs.AsAwsLogger(ctx)),
		config.WithRegion(envs.AwsRegion),
		config.WithEndpointResolver(customEndpointResolver),
	)
	if err != nil {
		logs.Logger(ctx).Errorf("unable to load SDK config, err=%v", err)
		return aws.Config{}, err
	}
	return cfg, nil
}
