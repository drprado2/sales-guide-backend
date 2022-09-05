package colors

import "github.com/mazznoer/csscolorparser"

func IsValid(color string) error {
	_, err := csscolorparser.Parse(color)
	return err
}

func AsRgb(color string) string {
	v, err := csscolorparser.Parse(color)
	if err != nil {
		return ""
	}
	return v.RGBString()
}

func AsHex(color string) string {
	v, err := csscolorparser.Parse(color)
	if err != nil {
		return ""
	}
	return v.HexString()
}
