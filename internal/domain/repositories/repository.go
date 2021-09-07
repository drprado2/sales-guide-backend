package repositories

import (
	"context"
	"github.com/drprado2/react-redux-typescript/internal/domain/entities"
)

type (
	CompanyRepository interface {
		Create(ctx context.Context, company *entities.Company) error

		GetCompanyByID(ctx context.Context, companyID string) (*entities.Company, error)
	}
)
