package utils

import "github.com/mazznoer/csscolorparser"

type (
	CssColorParserService struct{}
)

func (*CssColorParserService) IsValid(color string) error {
	_, err := csscolorparser.Parse(color)
	return err
}

func (*CssColorParserService) AsRgb(color string) string {
	v, err := csscolorparser.Parse(color)
	if err != nil {
		return ""
	}
	return v.RGBString()
}

func (*CssColorParserService) AsHex(color string) string {
	v, err := csscolorparser.Parse(color)
	if err != nil {
		return ""
	}
	return v.HexString()
}
