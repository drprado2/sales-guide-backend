package repositories

import (
	"context"
	"github.com/drprado2/react-redux-typescript/internal/domain/entities"
	"github.com/jackc/pgx/v4"
)

type (
	CompanyRepository interface {
		Create(ctx context.Context, company *entities.Company) (uint32, error)
		GetCompanyByID(ctx context.Context, companyID string) (*entities.Company, error)
	}

	UserRepository interface {
		Create(ctx context.Context, user *entities.User) (uint32, error)
		CreateTx(ctx context.Context, tx pgx.Tx, user *entities.User) (uint32, error)
		BeginTx(ctx context.Context) (pgx.Tx, error)
	}
)
