package logs

import (
	"context"
	"github.com/aws/smithy-go/logging"
	"github.com/drprado2/react-redux-typescript/configs"
	log "github.com/sirupsen/logrus"
	"os"
)

type (
	AwsLogger struct {
		logger *log.Entry
	}
)

func (al *AwsLogger) Logf(classification logging.Classification, format string, v ...interface{}) {
	switch classification {
	case logging.Warn:
		al.logger.Warnf(format, v)
	default:
		al.logger.Debug(format, v)
	}
}

func Setup() {
	envs := configs.Get()
	if envs.ServerEnvironment == configs.DeveloperEnvironment || envs.ServerEnvironment == configs.TestEnvironment {
		log.SetFormatter(&log.TextFormatter{
			FullTimestamp:             true,
			ForceColors:               true,
			EnvironmentOverrideColors: true,
		})
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetFormatter(&log.JSONFormatter{})
		log.SetLevel(log.InfoLevel)
	}

	log.SetOutput(os.Stdout)
	log.SetReportCaller(true)
}

func Logger(ctx context.Context) *log.Entry {
	contextLogger := log.WithFields(getDefaultLogFields(ctx))
	return contextLogger
}

func AsAwsLogger(ctx context.Context) logging.Logger {
	contextLogger := log.WithFields(getDefaultLogFields(ctx))
	return &AwsLogger{logger: contextLogger}
}

func getDefaultLogFields(ctx context.Context) log.Fields {
	envs := configs.Get()
	cid := ctx.Value("cid")
	if cid == nil {
		cid = "empty"
	}
	httpMethod := ctx.Value("httpMethod")
	httpPath := ctx.Value("httpPath")

	fields := log.Fields{
		"cid":     cid,
		"version": envs.SystemVersion,
		"app":     envs.AppName,
	}
	if httpMethod != nil {
		fields["httpMethod"] = httpMethod
	}
	if httpPath != nil {
		fields["httpPath"] = httpPath
	}
	return fields
}
