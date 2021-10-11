package usecases

import (
	"github.com/drprado2/react-redux-typescript/internal/domain/repositories"
	apptracer2 "github.com/drprado2/react-redux-typescript/pkg/apptracer"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"gopkg.in/auth0.v5/management"
)

var (
	companyRepository   repositories.CompanyRepository
	userRepository      repositories.UserRepository
	tracer              *apptracer2.TracerService
	payloadValidator    *validator.Validate
	validatorTranslates ut.Translator
	auth0manager        *management.Management
)

func Setup(
	companyRepositorySvc repositories.CompanyRepository,
	userRepositorySvc repositories.UserRepository,
	tracerSvc *apptracer2.TracerService,
	payloadValidatorSvc *validator.Validate,
	trans ut.Translator,
	auth0managerSvc *management.Management,
) {
	companyRepository = companyRepositorySvc
	userRepository = userRepositorySvc
	tracer = tracerSvc
	payloadValidator = payloadValidatorSvc
	validatorTranslates = trans
	auth0manager = auth0managerSvc
}
