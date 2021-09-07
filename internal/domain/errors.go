package domain

import (
	"errors"
	"fmt"
)

type (
	InternalError struct {
		err error
	}
)

var (
	CompanyInvalidToSaveError = errors.New("empresa inválida, por favor verifique os campos obrigatórios")
)

func NewInternalError(err error) error {
	return &InternalError{
		err: err,
	}
}

func (ie *InternalError) Error() string {
	return fmt.Sprintf("an expected error occur, err=%v", ie.err.Error())
}
