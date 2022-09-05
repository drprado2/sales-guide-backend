package validators

import (
	"errors"
	"github.com/paemuri/brdoc"
)

var (
	InvalidCnpjError = errors.New("O CNPJ é inválido")
)

func ValidateCnpj(document string) error {
	if !brdoc.IsCNPJ(document) {
		return InvalidCnpjError
	}
	return nil
}
