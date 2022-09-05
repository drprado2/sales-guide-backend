package domain

import (
	"github.com/drprado2/sales-guide/internal/domain/repositories"
	"github.com/drprado2/sales-guide/pkg/instrumentation/apptracer"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"gopkg.in/auth0.v5/management"
)

type (
	ServiceManager struct {
		CompanyRepository   repositories.CompanyRepository
		UserRepository      repositories.UserRepository
		Tracer              *apptracer.TracerService
		PayloadValidator    *validator.Validate
		ValidatorTranslates ut.Translator
		Auth0manager        *management.Management
	}
)

func CreateServiceManager(
	companyRepository repositories.CompanyRepository,
	userRepository repositories.UserRepository,
	tracer *apptracer.TracerService,
	payloadValidator *validator.Validate,
	validatorTranslator ut.Translator,
	auth0manager *management.Management,
) *ServiceManager {
	return &ServiceManager{
		CompanyRepository:   companyRepository,
		UserRepository:      userRepository,
		Tracer:              tracer,
		PayloadValidator:    payloadValidator,
		ValidatorTranslates: validatorTranslator,
		Auth0manager:        auth0manager,
	}
}
