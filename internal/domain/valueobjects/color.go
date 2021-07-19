package valueobjects

type (
	Color string

	ColorService interface {
		IsValid(color string) error
		AsRgb(color string) string
		AsHex(color string) string
	}
)

func NewColor(value string) (*Color, error) {
	if err := ColorSvc.IsValid(value); err != nil {
		return nil, err
	}
	return (*Color)(&value), nil
}

func (c *Color) Rgb() string {
	return ColorSvc.AsRgb(string(*c))
}

func (c *Color) Hex() string {
	return ColorSvc.AsHex(string(*c))
}
