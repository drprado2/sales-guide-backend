package configs

import (
	"fmt"
	"os"
	"strconv"
)

const (
	DeveloperEnvironment = "dev"
	TestEnvironment = "test"
	ProductionEnvironment = "prod"
)

type EnvsInterface interface {
	GetServerPort() int
	GetEnvironment() string
	GetDbUser() string
	GetDbPassword() string
	GetDbName() string
}

type Envs struct {}

func (*Envs) GetServerPort() int {
	v := os.Getenv("SERVER_PORT")
	if v != "" {
		vi, err := strconv.Atoi(v)
		if err != nil {
			panic(fmt.Sprintf("fail getting server port env, %v", err))
		}
		return vi
	}
	return 5050
}

func (*Envs) GetEnvironment() string {
	v := os.Getenv("SERVER_ENVIRONMENT")
	if v != "" {
		return v
	}
	return DeveloperEnvironment
}

func (*Envs) GetDbUser() string {
	v := os.Getenv("DB_USER")
	if v != "" {
		return v
	}
	return "postgres"
}

func (*Envs) GetDbPassword() string {
	v := os.Getenv("DB_PASS")
	if v != "" {
		return v
	}
	return "admin123"
}

func (*Envs) GetDbName() string {
	v := os.Getenv("DB_NAME")
	if v != "" {
		return v
	}
	return "poker-simulator"
}

func (*Envs) GetDbHost() string {
	v := os.Getenv("DB_HOST")
	if v != "" {
		return v
	}
	return "localhost"
}

func (*Envs) GetDbPort() string {
	v := os.Getenv("DB_PORT")
	if v != "" {
		return v
	}
	return "5432"
}

func (*Envs) GetZipkinUrl() string {
	v := os.Getenv("ZIPKIN_URL")
	if v != "" {
		return v
	}
	return "http://localhost:9411/api/v2/spans"
}
