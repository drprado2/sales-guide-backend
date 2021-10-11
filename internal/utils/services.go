package utils

type (
	ColorService interface {
		IsValid(color string) error
		AsRgb(color string) string
		AsHex(color string) string
	}

	CnpjValidator interface {
		Validate(document string) error
	}
)

var (
	ColorSvc         ColorService
	CnpjValidatorSvc CnpjValidator
)

func Setup(colorService ColorService, cnpjValidator CnpjValidator) {
	ColorSvc = colorService
	CnpjValidatorSvc = cnpjValidator
}
