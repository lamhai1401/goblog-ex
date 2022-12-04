package util

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validater linter
var Validater *Validator

// Validator linter
type Validator struct {
	validate *validator.Validate
}

func init() {
	// init validate
	validate := validator.New()
	// register function to get tag name from json tags.
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	Validater = &Validator{
		validate: validate,
	}
}

// ValidateStruct linter
func (s *Validator) ValidateStruct(data interface{}) error {
	return s.validate.Struct(s)
}
