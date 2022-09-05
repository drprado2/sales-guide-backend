package repository

import (
	"context"
	"github.com/drprado2/sales-guide/internal/domain/entities"
	apptracer2 "github.com/drprado2/sales-guide/pkg/instrumentation/apptracer"
	"github.com/drprado2/sales-guide/pkg/instrumentation/logs"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	InsertCompanyQuery = `
INSERT INTO company (
	id,
	name, 
	document, 
	logo, 
	total_colaborators, 
	primary_color, 
	primary_font_color, 
	secondary_color, 
	secondary_font_color, 
	created_at, 
	updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	RETURNING xmin
`

	GetCompanyByIdQuery = `
SELECT
	id,
	name, 
	document, 
	logo, 
	total_colaborators, 
	primary_color, 
	primary_font_color, 
	secondary_color, 
	secondary_font_color, 
	created_at, 
	updated_at,
	xmin
FROM company
WHERE id = $1
`
)

type (
	CompanySqlRepository struct {
		dbPool    *pgxpool.Pool
		tracerSvc *apptracer2.TracerService
	}
)

func NewCompanySqlRepository(dbPool *pgxpool.Pool, tracerSvc *apptracer2.TracerService) *CompanySqlRepository {
	return &CompanySqlRepository{
		dbPool:    dbPool,
		tracerSvc: tracerSvc,
	}
}

func (csr *CompanySqlRepository) Create(ctx context.Context, company *entities.Company) (uint32, error) {
	span, ctx := csr.tracerSvc.SpanFromContext(ctx)
	defer span.Finish()

	var rowversion uint32
	row := csr.dbPool.QueryRow(
		ctx,
		InsertCompanyQuery,
		company.ID,
		company.Name,
		company.Document.AsString(),
		company.Logo.AsString(),
		company.TotalColaborators,
		company.PrimaryColor.Hex(),
		company.PrimaryFontColor.Hex(),
		company.SecondaryColor.Hex(),
		company.SecondaryFontColor.Hex(),
		company.CreatedAt,
		company.UpdatedAt,
	)
	err := row.Scan(&rowversion)
	return rowversion, err
}

func (csr *CompanySqlRepository) GetCompanyByID(ctx context.Context, companyID string) (*entities.Company, error) {
	span, ctx := csr.tracerSvc.SpanFromContext(ctx)
	defer span.Finish()

	result := &entities.Company{}
	if err := csr.dbPool.QueryRow(ctx, GetCompanyByIdQuery, companyID).
		Scan(&result.ID, &result.Name, &result.Document, &result.Logo, &result.TotalColaborators, &result.PrimaryColor, &result.PrimaryFontColor, &result.SecondaryColor, &result.SecondaryFontColor, &result.CreatedAt, &result.UpdatedAt, &result.RowVersion); err != nil {
		logs.Logger(ctx).WithError(err).Error("Error getting player by id")
		return nil, err
	}
	return result, nil
}
