package configs

import (
	"log"
	"sync"

	"github.com/gosidekick/goconfig"
)

const (
	DeveloperEnvironment  = "dev"
	TestEnvironment       = "test"
	ProductionEnvironment = "prod"
)

var (
	doOnce sync.Once
	env    *Environment
)

//Environment this object keep the all variables environment
type Environment struct {
	ServerPort         int    `cfg:"SERVER_PORT" cfgDefault:"5050" cfgRequired:"true"`
	ServerEnvironment  string `cfg:"SERVER_ENVIRONMENT" cfgDefault:"dev" cfgRequired:"true"`
	SystemVersion      string `cfg:"SYSTEM_VERSION" cfgDefault:"UNKNOWN" cfgRequired:"true"`
	AppName            string `cfg:"APP_NAME" cfgDefault:"api-sales-guide" cfgRequired:"true"`
	DbUser             string `cfg:"DB_USER" cfgDefault:"postgres" cfgRequired:"true"`
	DbPass             string `cfg:"DB_PASS" cfgDefault:"admin123" cfgRequired:"true"`
	DbName             string `cfg:"DB_NAME" cfgDefault:"sales-guide" cfgRequired:"true"`
	DbHost             string `cfg:"DB_HOST" cfgDefault:"localhost" cfgRequired:"true"`
	DbPort             int    `cfg:"DB_PORT" cfgDefault:"4611" cfgRequired:"true"`
	ZipkinUrl          string `cfg:"ZIPKIN_URL" cfgDefault:"http://localhost:9411/api/v2/spans" cfgRequired:"true"`
	AwsRegion          string `cfg:"AWS_REGION" cfgDefault:"sa-east-1" cfgRequired:"true"`
	AwsAccessKey       string `cfg:"AWS_ACCESS_KEY_ID" cfgDefault:"sa-east-1" cfgRequired:"true"`
	AwsSecretAccessKey string `cfg:"AWS_SECRET_ACCESS_KEY" cfgDefault:"sa-east-1" cfgRequired:"true"`
	AwsEndpoint        string `cfg:"AWS_ENDPOINT" cfgDefault:"http://localhost:4566" cfgRequired:"true"`
	RedisHost          string `cfg:"REDIS_HOST" cfgDefault:"localhost" cfgRequired:"true"`
	RedisPort          int    `cfg:"REDIS_PORT" cfgDefault:"4511" cfgRequired:"true"`
	RedisPass          string `cfg:"REDIS_PASS" cfgDefault:"" cfgRequired:"false"`
	RedisDb            int    `cfg:"REDIS_DB" cfgDefault:"0" cfgRequired:"false"`
	ForceS3PathStyle   bool   `cfg:"FORCE_S3_PATH_STYLE" cfgDefault:"true" cfgRequired:"false"`
}

//Get return the instance of environment that keep the environment variables
func Get() *Environment {
	doOnce.Do(func() {
		env = &Environment{}
		err := goconfig.Parse(env)
		if err != nil {
			log.Fatal(err)
		}
	})
	return env
}

//Reset will reload the environment variables
func Reset() *Environment {
	doOnce = sync.Once{}
	return Get()
}
