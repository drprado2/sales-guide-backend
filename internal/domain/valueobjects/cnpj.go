package valueobjects

import (
	"github.com/drprado2/sales-guide/pkg/validators"
)

type (
	Cnpj string
)

func (c *Cnpj) AsString() string {
	return (string)(*c)
}

func NewCnpj(document string) (*Cnpj, error) {
	if err := validators.ValidateCnpj(document); err != nil {
		return nil, err
	}
	return (*Cnpj)(&document), nil
}
