package repository

import (
	"context"
	"fmt"
	"github.com/drprado2/sales-guide/configs"
	"github.com/drprado2/sales-guide/pkg/instrumentation/logs"
	"github.com/jackc/pgx/v4/pgxpool"
)

func CreateConnPool(ctx context.Context, config *configs.Config) (*pgxpool.Pool, error) {
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%v/%s?connect_timeout=15&application_name=sales-guide",
		config.DbUser,
		config.DbPass,
		config.DbHost,
		config.DbPort,
		config.DbName)
	logs.Logger(ctx).Info("Connecting postgres DB")
	dbpool, err := pgxpool.Connect(ctx, connectionString)
	if err != nil {
		logs.Logger(ctx).Errorf("error creating DB connection, %v", err)
		return nil, err
	}
	logs.Logger(ctx).Info("DB connected successfully")
	return dbpool, err
}
