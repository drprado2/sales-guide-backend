package repository

import (
	"context"
	"github.com/drprado2/sales-guide/internal/domain/entities"
	apptracer2 "github.com/drprado2/sales-guide/pkg/instrumentation/apptracer"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	InsertUserQuery = `
INSERT INTO app_user (
	id,
	company_id,
	name, 
	email, 
	phone, 
	birth_date, 
	avatar_image, 
	record_creation_count, 
	record_editing_count, 
	record_deletion_count, 
	created_at, 
	updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	RETURNING xmin
`
)

type (
	UserSqlRepository struct {
		dbPool    *pgxpool.Pool
		tracerSvc *apptracer2.TracerService
	}
)

func NewUserSqlRepository(dbPool *pgxpool.Pool, tracerSvc *apptracer2.TracerService) *UserSqlRepository {
	return &UserSqlRepository{
		dbPool:    dbPool,
		tracerSvc: tracerSvc,
	}
}

func (csr *UserSqlRepository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	span, ctx := csr.tracerSvc.SpanFromContext(ctx)
	defer span.Finish()

	return csr.dbPool.Begin(ctx)
}

func (csr *UserSqlRepository) Create(ctx context.Context, user *entities.User) (uint32, error) {
	span, ctx := csr.tracerSvc.SpanFromContext(ctx)
	defer span.Finish()

	var rowversion uint32
	err := csr.dbPool.QueryRow(
		ctx,
		InsertUserQuery,
		user.ID,
		user.CompanyID,
		user.Name,
		user.Email,
		user.Phone,
		user.BirthDate,
		user.AvatarImage.AsString(),
		user.RecordCreationCount,
		user.RecordEditingCount,
		user.RecordDeletionCount,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&rowversion)

	return rowversion, err
}

func (csr *UserSqlRepository) CreateTx(ctx context.Context, tx pgx.Tx, user *entities.User) (uint32, error) {
	span, ctx := csr.tracerSvc.SpanFromContext(ctx)
	defer span.Finish()

	var rowversion uint32
	err := tx.QueryRow(
		ctx,
		InsertUserQuery,
		user.ID,
		user.CompanyID,
		user.Name,
		user.Email,
		user.Phone,
		user.BirthDate,
		user.AvatarImage.AsString(),
		user.RecordCreationCount,
		user.RecordEditingCount,
		user.RecordDeletionCount,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&rowversion)
	return rowversion, err
}
