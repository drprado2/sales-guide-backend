package validations

import (
	"github.com/go-playground/locales/pt_BR"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	pt_translations "github.com/go-playground/validator/v10/translations/pt_BR"
	"reflect"
	"strings"
)

func CreateValidatorService() (*validator.Validate, ut.Translator) {
	ptbr := pt_BR.New()
	uni := ut.New(ptbr, ptbr)
	trans, _ := uni.GetTranslator("pt_BR")
	v := validator.New()
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("name"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	if err := pt_translations.RegisterDefaultTranslations(v, trans); err != nil {
		panic(err)
	}
	return v, trans
}
