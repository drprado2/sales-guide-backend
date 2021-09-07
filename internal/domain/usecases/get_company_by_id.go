package usecases

import (
	"context"
	"time"
)

type (
	GetCompanyOutput struct {
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

func GetCompanyByID(ctx context.Context, id string) (*GetCompanyOutput, error) {
	company, err := companyRepository.GetCompanyByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if company == nil {
		return nil, nil
	}
	return &GetCompanyOutput{
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
