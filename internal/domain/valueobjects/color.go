package valueobjects

import (
	"github.com/drprado2/sales-guide/pkg/colors"
)

type (
	Color string
)

func NewColor(value string) (*Color, error) {
	if err := colors.IsValid(value); err != nil {
		return nil, err
	}
	return (*Color)(&value), nil
}

func (c *Color) Rgb() string {
	return colors.AsRgb(string(*c))
}

func (c *Color) Hex() string {
	return colors.AsHex(string(*c))
}
