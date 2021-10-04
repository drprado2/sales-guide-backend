package usecases

import (
	"context"
	"errors"
	"github.com/drprado2/react-redux-typescript/internal/domain"
	"github.com/drprado2/react-redux-typescript/internal/domain/entities"
	"github.com/drprado2/react-redux-typescript/internal/domain/valueobjects"
	"github.com/drprado2/react-redux-typescript/pkg/logs"
	"github.com/google/uuid"
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
	}
)

func CreateCompany(ctx context.Context, input *CreateCompanyInput) (*CreateCompanyOutput, error) {
	span, ctx := tracer.SpanFromContext(ctx)
	defer span.Finish()

	if err := payloadValidator.Struct(input); err != nil {
		return nil, domain.PayloadErrorFromValidator(err, validatorTranslates)
	}

	errGroup := &domain.ErrorGroup{}

	cnpj, err := valueobjects.NewCnpj(input.Document)
	if err != nil {
		logs.Logger(ctx).WithError(err).Warn("fail creating company with invalid CNPJ")
		errGroup.Append(domain.NewConstraintError(err))
	}
	company := entities.NewCompany(uuid.NewString(), input.Name, cnpj)
	if input.Logo != nil {
		logoUri, err := valueobjects.NewUri(*input.Logo)
		if err != nil {
			logs.Logger(ctx).WithError(err).Warn("fail creating company with invalid logo URI")
			errGroup.Append(domain.NewConstraintError(errors.New("logo marca inválida")))
		}
		company.Logo = logoUri
	}
	if input.PrimaryColor != nil {
		color, err := valueobjects.NewColor(*input.PrimaryColor)
		if err != nil {
			logs.Logger(ctx).WithError(err).Warn("fail creating company with invalid primary color")
			errGroup.Append(domain.NewConstraintError(errors.New("cor primária inválida")))
		}
		company.PrimaryColor = color
	}
	if input.PrimaryFontColor != nil {
		color, err := valueobjects.NewColor(*input.PrimaryFontColor)
		if err != nil {
			logs.Logger(ctx).WithError(err).Warn("fail creating company with invalid primary font color")
			errGroup.Append(domain.NewConstraintError(errors.New("cor de fonte primária inválida")))
		}
		company.PrimaryFontColor = color
	}
	if input.SecondaryColor != nil {
		color, err := valueobjects.NewColor(*input.SecondaryColor)
		if err != nil {
			logs.Logger(ctx).WithError(err).Warn("fail creating company with invalid secondary color")
			errGroup.Append(domain.NewConstraintError(errors.New("cor secundária inválida")))
		}
		company.SecondaryColor = color
	}
	if input.SecondaryFontColor != nil {
		color, err := valueobjects.NewColor(*input.SecondaryFontColor)
		if err != nil {
			logs.Logger(ctx).WithError(err).Warn("fail creating company with invalid secondary font color")
			errGroup.Append(domain.NewConstraintError(errors.New("cor de fonte secundária inválida")))
		}
		company.SecondaryFontColor = color
	}

	if errGroup.HasError() {
		return nil, errGroup.AsError()
	}

	company.TotalColaborators = 0

	if err := companyRepository.Create(ctx, company); err != nil {
		logs.Logger(ctx).WithError(err).Error("fail creating company with error in repository")
		return nil, err
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
		CreatedAt:          company.CreatedAt,
		UpdatedAt:          company.UpdatedAt,
	}, nil
}
