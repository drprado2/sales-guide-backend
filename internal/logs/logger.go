package logs

import (
	"context"
	"github.com/drprado2/react-redux-typescript/configs"
	log "github.com/sirupsen/logrus"
	"os"
)

func Init(envs configs.EnvsInterface) {
	if envs.GetEnvironment() == configs.DeveloperEnvironment || envs.GetEnvironment() == configs.TestEnvironment {
		log.SetFormatter(&log.TextFormatter{
			FullTimestamp: true,
			ForceColors:   true,
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

func getDefaultLogFields(ctx context.Context) log.Fields {
	cid := ctx.Value("cid")
	if cid == nil {
		cid = "empty"
	}
	httpMethod := ctx.Value("httpMethod")
	httpPath := ctx.Value("httpPath")

	fields := log.Fields{
		"cid": cid,
	}
	if httpMethod != nil {
		fields["httpMethod"] = httpMethod
	}
	if httpPath != nil {
		fields["httpPath"] = httpPath
	}
	return fields
}

