package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/drprado2/sales-guide/configs"
	"github.com/drprado2/sales-guide/pkg/instrumentation/logs"
	"github.com/go-redis/redis/v8"
	"time"
)

var (
	client *redis.Client
)

func Setup(ctx context.Context) {
	logs.Logger(ctx).Infof("Starting Redis setup")
	envs := configs.Get()

	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%v", envs.RedisHost, envs.RedisPort),
		Password: envs.RedisPass,
		DB:       envs.RedisDb,
	})
	logs.Logger(ctx).Infof("Redis setup fineshed")
}

func GetFromKeySvc(ctx context.Context, key string) (*string, error) {
	value, err := client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &value, nil
}

func UnmarshalFromKeySvc(ctx context.Context, key string, target interface{}) (bool, error) {
	value, err := client.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if err := json.Unmarshal([]byte(value), target); err != nil {
		return false, err
	}
	return true, nil
}

func SetKeySvc(ctx context.Context, key string, value string, ttl time.Duration) error {
	return client.Set(ctx, key, value, ttl).Err()
}

func SetKeyJsonSvc(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	val, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return client.Set(ctx, key, val, ttl).Err()
}

func DeleteKeySvc(ctx context.Context, key string) error {
	return client.Del(ctx, key).Err()
}
