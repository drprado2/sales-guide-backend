package ses

import (
	"errors"
	"strings"
)

var (
	RequiredSenderError       = errors.New("sender is required")
	RequiredCharsetError      = errors.New("charset is required")
	RequiredContentError      = errors.New("content HTML or text is required")
	RequiredSubjectError      = errors.New("subject is required")
	RequiredTemplateError     = errors.New("template name and template data are required")
	RequiredDestinationsError = errors.New("you must to inform at least one destination")
)

func GroupErrors(errs ...error) error {
	errsText := make([]string, len(errs))
	for i, err := range errs {
		errsText[i] = err.Error()
	}
	return errors.New(strings.Join(errsText, ", "))
}
