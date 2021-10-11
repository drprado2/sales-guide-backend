package domain

import (
	"errors"
	"fmt"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"strings"
)

const (
	internalErrorType           = "InternalError"
	payloadErrorType            = "PayloadError"
	constraintErrorType         = "ConstraintError"
	PgUniqueConstraintErrorCode = "23505"
	PgForeignErrorCode          = "23503"
)

type (
	InternalError struct {
		err error
	}
	PayloadError struct {
		err error
	}
	ConstraintError struct {
		err error
	}
	ErrorGroup struct {
		errors []error
	}
)

var (
	CompanyInvalidToSaveError = errors.New("empresa inválida, por favor verifique os campos obrigatórios")
	UserInvalidToSaveError    = errors.New("usuário inválido, por favor verifique os campos obrigatórios")
)

func NewInternalError(err error) error {
	return &InternalError{
		err: err,
	}
}

func NewPayloadError(err error) error {
	return &PayloadError{
		err: err,
	}
}

func PayloadErrorFromValidator(err error, trans ut.Translator) error {
	errs := err.(validator.ValidationErrors)
	sb := strings.Builder{}
	for _, e := range errs {
		sb.WriteString(e.Translate(trans) + ", ")
	}
	finalErr := sb.String()
	return &PayloadError{
		err: errors.New(finalErr[0 : len(finalErr)-2]),
	}
}

func NewConstraintError(err error) error {
	return &ConstraintError{
		err: err,
	}
}

func (ie *InternalError) Error() string {
	return fmt.Sprintf("an expected error occur, err=%v", ie.err.Error())
}

func (pe *PayloadError) Error() string {
	return fmt.Sprintf("%v", pe.err)
}

func (ce *ConstraintError) Error() string {
	return fmt.Sprintf("%v", ce.err.Error())
}

func (ce *ErrorGroup) Append(err error) {
	ce.errors = append(ce.errors, err)
}

func (ce *ErrorGroup) HasError() bool {
	return len(ce.errors) > 0
}

func (ce *ErrorGroup) AsError() error {
	sb := strings.Builder{}
	errorType := payloadErrorType
	for i, e := range ce.errors {
		if i > 0 {
			sb.WriteString(fmt.Sprintf(", %s", e.Error()))
		} else {
			sb.WriteString(e.Error())
		}

		if _, ok := e.(*InternalError); ok {
			errorType = internalErrorType
		} else if _, ok := e.(*ConstraintError); ok && errorType != internalErrorType {
			errorType = constraintErrorType
		}
	}
	switch errorType {
	case internalErrorType:
		return NewInternalError(errors.New(sb.String()))
	case constraintErrorType:
		return NewConstraintError(errors.New(sb.String()))
	default:
		return NewPayloadError(errors.New(sb.String()))
	}
}
