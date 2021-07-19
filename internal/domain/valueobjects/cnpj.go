package valueobjects

type (
	Cnpj string

	CnpjValidator interface {
		Validate(document string) error
	}
)

func (c *Cnpj) AsString() string {
	return (string)(*c)
}

func NewCnpj(document string) (*Cnpj, error) {
	if err := CnpjValidatorSvc.Validate(document); err != nil {
		return nil, err
	}
	return (*Cnpj)(&document), nil
}
