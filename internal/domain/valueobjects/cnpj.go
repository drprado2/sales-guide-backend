package valueobjects

import (
	utils2 "github.com/drprado2/react-redux-typescript/internal/utils"
)

type (
	Cnpj string
)

func (c *Cnpj) AsString() string {
	return (string)(*c)
}

func NewCnpj(document string) (*Cnpj, error) {
	if err := utils2.CnpjValidatorSvc.Validate(document); err != nil {
		return nil, err
	}
	return (*Cnpj)(&document), nil
}
