package usecases

import (
	"github.com/drprado2/react-redux-typescript/internal/domain/repositories"
	apptracer2 "github.com/drprado2/react-redux-typescript/pkg/apptracer"
)

var (
	companyRepository repositories.CompanyRepository
	tracer            *apptracer2.TracerService
)

func Setup(
	companyRepositorySvc repositories.CompanyRepository,
	tracerSvc *apptracer2.TracerService,
) {
	companyRepository = companyRepositorySvc
	tracer = tracerSvc
}
