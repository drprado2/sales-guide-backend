package usecases

import (
	"context"
	"errors"
	"github.com/drprado2/sales-guide/internal/domain"
	"github.com/drprado2/sales-guide/internal/domain/entities"
	errors2 "github.com/drprado2/sales-guide/internal/domain/errors"
	"github.com/drprado2/sales-guide/internal/domain/valueobjects"
	"github.com/drprado2/sales-guide/pkg/ctxvals"
	"github.com/drprado2/sales-guide/pkg/instrumentation/logs"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"time"
)

type (
	CreateCompanyInput struct {
		Name               string  `json:"name,omitempty" validate:"required,min=1" name:"nome"`
		Document           string  `json:"document,omitempty" validate:"required,min=14,max=14" name:"documento"`
		Logo               *string `json:"logo,omitempty"`
		PrimaryColor       *string `json:"primaryColor,omitempty"`
		PrimaryFontColor   *string `json:"primaryFontColor,omitempty"`
		SecondaryColor     *string `json:"secondaryColor,omitempty"`
		SecondaryFontColor *string `json:"secondaryFontColor,omitempty"`
	}

	CreateCompanyOutput struct {
		Id                 string    `json:"id,omitempty"`
		Name               string    `json:"name,omitempty"`
		Document           string    `json:"document,omitempty"`
		Logo               string    `json:"logo,omitempty"`
		TotalColaborators  int       `json:"totalColaborators,omitempty"`
		PrimaryColor       string    `json:"primaryColor,omitempty"`
		PrimaryFontColor   string    `json:"primaryFontColor,omitempty"`
		SecondaryColor     string    `json:"secondaryColor,omitempty"`
		SecondaryFontColor string    `json:"secondaryFontColor,omitempty"`
		CreatedAt          time.Time `json:"createdAt,omitempty"`
		UpdatedAt          time.Time `json:"UpdatedAt,omitempty"`
		RowVersion         uint32    `json:"rowVersion,omitempty"`
	}
)

func CreateCompany(ctx context.Context, sm *domain.ServiceManager, input *CreateCompanyInput) (*CreateCompanyOutput, error) {
	span, ctx := sm.Tracer.SpanFromContext(ctx)
	defer span.Finish()

	location := ctxvals.LocationOrDefault(ctx)

	if err := sm.PayloadValidator.Struct(input); err != nil {
		return nil, errors2.PayloadErrorFromValidator(err, sm.ValidatorTranslates)
	}

	errGroup := &errors2.ErrorGroup{}

	cnpj, err := valueobjects.NewCnpj(input.Document)
	if err != nil {
		logs.Logger(ctx).WithError(err).Warn("fail creating company with invalid CNPJ")
		errGroup.Append(errors2.NewConstraintError(err))
	}
	company := entities.NewCompany(uuid.NewString(), input.Name, cnpj)
	if input.Logo != nil {
		logoUri, err := valueobjects.NewUri(*input.Logo)
		if err != nil {
			logs.Logger(ctx).WithError(err).Warn("fail creating company with invalid logo URI")
			errGroup.Append(errors2.NewConstraintError(errors.New("logo marca inválida")))
		}
		company.Logo = logoUri
	}
	if input.PrimaryColor != nil {
		color, err := valueobjects.NewColor(*input.PrimaryColor)
		if err != nil {
			logs.Logger(ctx).WithError(err).Warn("fail creating company with invalid primary color")
			errGroup.Append(errors2.NewConstraintError(errors.New("cor primária inválida")))
		}
		company.PrimaryColor = color
	}
	if input.PrimaryFontColor != nil {
		color, err := valueobjects.NewColor(*input.PrimaryFontColor)
		if err != nil {
			logs.Logger(ctx).WithError(err).Warn("fail creating company with invalid primary font color")
			errGroup.Append(errors2.NewConstraintError(errors.New("cor de fonte primária inválida")))
		}
		company.PrimaryFontColor = color
	}
	if input.SecondaryColor != nil {
		color, err := valueobjects.NewColor(*input.SecondaryColor)
		if err != nil {
			logs.Logger(ctx).WithError(err).Warn("fail creating company with invalid secondary color")
			errGroup.Append(errors2.NewConstraintError(errors.New("cor secundária inválida")))
		}
		company.SecondaryColor = color
	}
	if input.SecondaryFontColor != nil {
		color, err := valueobjects.NewColor(*input.SecondaryFontColor)
		if err != nil {
			logs.Logger(ctx).WithError(err).Warn("fail creating company with invalid secondary font color")
			errGroup.Append(errors2.NewConstraintError(errors.New("cor de fonte secundária inválida")))
		}
		company.SecondaryFontColor = color
	}
	if err := company.Validate(); err != nil {
		logs.Logger(ctx).WithError(err).Warn("company invalid to save")
		errGroup.Append(errors2.NewConstraintError(err))
	}

	if errGroup.HasError() {
		return nil, errGroup.AsError()
	}

	company.TotalColaborators = 0

	rowversion, err := sm.CompanyRepository.Create(ctx, company)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			if err.Code == errors2.PgUniqueConstraintErrorCode {
				logs.Logger(ctx).WithError(err).Error("fail creating company with unique constraint error")
				return nil, errors2.NewConstraintError(errors.New("já existe uma empresa com esse documento"))
			}
		}
		logs.Logger(ctx).WithError(err).Error("fail creating company with error in repository")
		return nil, errors2.NewInternalError(err)
	}

	return &CreateCompanyOutput{
		Id:                 company.ID,
		Name:               company.Name,
		Document:           company.Document.AsString(),
		Logo:               company.Logo.AsString(),
		TotalColaborators:  company.TotalColaborators,
		PrimaryColor:       company.PrimaryColor.Hex(),
		PrimaryFontColor:   company.PrimaryFontColor.Hex(),
		SecondaryColor:     company.SecondaryColor.Hex(),
		SecondaryFontColor: company.SecondaryFontColor.Hex(),
		CreatedAt:          company.CreatedAt.In(location),
		UpdatedAt:          company.UpdatedAt.In(location),
		RowVersion:         rowversion,
	}, nil
}
