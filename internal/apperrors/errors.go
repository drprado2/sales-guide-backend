package apperrors

import "errors"

var (
	InvalidRequestParameters = errors.New("Your request has invalid paraters")
)
