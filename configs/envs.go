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
	ServerPort            int    `cfg:"SERVER_PORT" cfgDefault:"5050" cfgRequired:"true"`
	ServerHost            string `cfg:"SERVER_HOST" cfgDefault:"localhost" cfgRequired:"true"`
	ServiceName           string `cfg:"SERVICE_NAME" cfgDefault:"Sales Guide" cfgRequired:"true"`
	ServerEndpointTimeout int    `cfg:"SERVER_ENDPOINT_TIMEOUT" cfgDefault:"90000" cfgRequired:"true"`
	ServerEnvironment     string `cfg:"SERVER_ENVIRONMENT" cfgDefault:"dev" cfgRequired:"true"`
	SystemVersion         string `cfg:"SYSTEM_VERSION" cfgDefault:"UNKNOWN" cfgRequired:"true"`
	AppName               string `cfg:"APP_NAME" cfgDefault:"api-sales-guide" cfgRequired:"true"`
	DbUser                string `cfg:"DB_USER" cfgDefault:"admin" cfgRequired:"true"`
	DbPass                string `cfg:"DB_PASS" cfgDefault:"Postgres2019!" cfgRequired:"true"`
	DbName                string `cfg:"DB_NAME" cfgDefault:"sales-guide" cfgRequired:"true"`
	DbHost                string `cfg:"DB_HOST" cfgDefault:"localhost" cfgRequired:"true"`
	DbPort                int    `cfg:"DB_PORT" cfgDefault:"5432" cfgRequired:"true"`
	ZipkinReportUrl       string `cfg:"ZIPKIN_URL" cfgDefault:"http://localhost:9411/api/v2/spans" cfgRequired:"true"`
	AwsRegion             string `cfg:"AWS_REGION" cfgDefault:"sa-east-1" cfgRequired:"true"`
	AwsAccessKey          string `cfg:"AWS_ACCESS_KEY_ID" cfgDefault:"AKIA4TOL2AU6VBTLRVV7" cfgRequired:"true"`
	AwsSecretAccessKey    string `cfg:"AWS_SECRET_ACCESS_KEY" cfgDefault:"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" cfgRequired:"true"`
	//AwsEndpoint        string `cfg:"AWS_ENDPOINT" cfgDefault:"http://localhost:4566" cfgRequired:"true"`
	AwsEndpoint       string `cfg:"AWS_ENDPOINT" cfgDefault:"http://d6be671397b8.ngrok.io:80" cfgRequired:"true"`
	RedisHost         string `cfg:"REDIS_HOST" cfgDefault:"localhost" cfgRequired:"true"`
	RedisPort         int    `cfg:"REDIS_PORT" cfgDefault:"4511" cfgRequired:"true"`
	RedisPass         string `cfg:"REDIS_PASS" cfgDefault:"" cfgRequired:"false"`
	RedisDb           int    `cfg:"REDIS_DB" cfgDefault:"0" cfgRequired:"false"`
	ForceS3PathStyle  bool   `cfg:"FORCE_S3_PATH_STYLE" cfgDefault:"true" cfgRequired:"false"`
	Auth0Domain       string `cfg:"AUTH0_DOMAIN" cfgDefault:"drprado2.us.auth0.com" cfgRequired:"true"`
	Auth0ClientID     string `cfg:"AUTH0_CLIENT_ID" cfgDefault:"qnDkRCLMB8MIcuFVaswOrJtG0aR1vpsy" cfgRequired:"true"`
	Auth0ClientSecret string `cfg:"AUTH0_CLIENT_SECRET" cfgDefault:"eujiEFfnph_GeJImDdIxFo9m2WZ9xOObNFiMXhLpuACJGPHPx_NlSVK4FocmI-AS" cfgRequired:"true"`
	Auth0VerifyEmail  bool   `cfg:"AUTH0_VERIFY_EMAIL" cfgDefault:"false" cfgRequired:"true"`
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
