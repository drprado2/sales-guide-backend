package valueobjects

var (
	ColorSvc         ColorService
	CnpjValidatorSvc CnpjValidator
)

func Setup(
	colorSvc ColorService,
	cnpjValidator CnpjValidator,
) {
	ColorSvc = colorSvc
	CnpjValidatorSvc = cnpjValidator
}
