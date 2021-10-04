package usecases

import (
	"github.com/drprado2/react-redux-typescript/internal/domain/repositories"
	apptracer2 "github.com/drprado2/react-redux-typescript/pkg/apptracer"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

var (
	companyRepository   repositories.CompanyRepository
	tracer              *apptracer2.TracerService
	payloadValidator    *validator.Validate
	validatorTranslates ut.Translator
)

func Setup(
	companyRepositorySvc repositories.CompanyRepository,
	tracerSvc *apptracer2.TracerService,
	payloadValidatorSvc *validator.Validate,
	trans ut.Translator,
) {
	companyRepository = companyRepositorySvc
	tracer = tracerSvc
	payloadValidator = payloadValidatorSvc
	validatorTranslates = trans
}
