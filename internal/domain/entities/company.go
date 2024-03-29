package entities

import (
	"github.com/drprado2/sales-guide/internal/domain/errors"

	//"github.com/drprado2/sales-guide/internal/domain"
	"github.com/drprado2/sales-guide/internal/domain/valueobjects"
	"strings"
	"time"
)

const (
	DefaultLogo               = "http://www.hsevolutione.com/wp-content/uploads/2020/07/logo-principal.png"
	DefaultPrimaryColor       = "#000066"
	DefaultPrimaryFontColor   = "#cce4ff"
	DefaultSecondaryColor     = "#ffffff"
	DefaultSecondaryFontColor = "#222"
)

type (
	Company struct {
		ID                 string
		Name               string
		Document           *valueobjects.Cnpj
		Logo               *valueobjects.Uri
		TotalColaborators  int
		PrimaryColor       *valueobjects.Color
		PrimaryFontColor   *valueobjects.Color
		SecondaryColor     *valueobjects.Color
		SecondaryFontColor *valueobjects.Color
		CreatedAt          time.Time
		UpdatedAt          time.Time
		RowVersion         uint32
	}
)

func NewCompany(id string, name string, cnpj *valueobjects.Cnpj) *Company {
	logoUri, _ := valueobjects.NewUri(DefaultLogo)
	pColor, _ := valueobjects.NewColor(DefaultPrimaryColor)
	pfColor, _ := valueobjects.NewColor(DefaultPrimaryFontColor)
	sColor, _ := valueobjects.NewColor(DefaultSecondaryColor)
	sfColor, _ := valueobjects.NewColor(DefaultSecondaryFontColor)

	return &Company{
		ID:                 id,
		Document:           cnpj,
		Name:               strings.ToLower(name),
		Logo:               logoUri,
		TotalColaborators:  0,
		PrimaryColor:       pColor,
		PrimaryFontColor:   pfColor,
		SecondaryColor:     sColor,
		SecondaryFontColor: sfColor,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
}

func (c *Company) Validate() error {
	if c.ID == "" ||
		c.Name == "" ||
		c.Logo == nil ||
		c.Document == nil ||
		c.PrimaryColor == nil ||
		c.PrimaryFontColor == nil ||
		c.SecondaryColor == nil ||
		c.SecondaryFontColor == nil {
		return errors.CompanyInvalidToSaveError
	}
	return nil
}
