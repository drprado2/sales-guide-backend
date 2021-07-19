package domain

import "fmt"

type (
	InternalError struct {
		err error
	}
)

func NewInternalError(err error) error {
	return &InternalError{
		err: err,
	}
}

func (ie *InternalError) Error() string {
	return fmt.Sprintf("an expected error occur, err=%v", ie.err.Error())
}
