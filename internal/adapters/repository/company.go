package repository

import (
	"context"
	"github.com/drprado2/react-redux-typescript/internal/domain/entities"
	apptracer2 "github.com/drprado2/react-redux-typescript/pkg/apptracer"
	"github.com/drprado2/react-redux-typescript/pkg/logs"
	"github.com/jackc/pgx/v4/pgxpool"
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

func (csr *CompanySqlRepository) Create(ctx context.Context, company *entities.Company) error {
	span, ctx := csr.tracerSvc.SpanFromContext(ctx)
	defer span.Finish()

	query := `
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
`
	_, err := csr.dbPool.Exec(
		ctx,
		query,
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
	return err
}

func (csr *CompanySqlRepository) GetCompanyByID(ctx context.Context, companyID string) (*entities.Company, error) {
	span, ctx := csr.tracerSvc.SpanFromContext(ctx)
	defer span.Finish()

	query := `
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
	updated_at
FROM company
WHERE id = $1
`
	result := &entities.Company{}
	if err := csr.dbPool.QueryRow(ctx, query, companyID).
		Scan(&result.ID, &result.Name, &result.Document, &result.Logo, &result.TotalColaborators, &result.PrimaryColor, &result.PrimaryFontColor, &result.SecondaryColor, &result.SecondaryFontColor, &result.CreatedAt, &result.UpdatedAt); err != nil {
			logs.Logger(ctx).WithError(err).Error("Error getting player by id")
		return nil, err
	}
	return result, nil
}
