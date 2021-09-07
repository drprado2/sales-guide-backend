package valueobjects

import (
	utils2 "github.com/drprado2/react-redux-typescript/internal/utils"
)

type (
	Color string
)

func NewColor(value string) (*Color, error) {
	if err := utils2.ColorSvc.IsValid(value); err != nil {
		return nil, err
	}
	return (*Color)(&value), nil
}

func (c *Color) Rgb() string {
	return utils2.ColorSvc.AsRgb(string(*c))
}

func (c *Color) Hex() string {
	return utils2.ColorSvc.AsHex(string(*c))
}
