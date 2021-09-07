package utils

import (
	"errors"
	"github.com/paemuri/brdoc"
)

type (
	PaemureBrDocCnpjValidator struct{}
)

var (
	InvalidCnpjError = errors.New("O CNPJ é inválido")
)

func (*PaemureBrDocCnpjValidator) Validate(document string) error {
	if !brdoc.IsCNPJ(document) {
		return InvalidCnpjError
	}
	return nil
}
