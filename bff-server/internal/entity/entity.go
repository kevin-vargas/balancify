package entity

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type User struct {
	Email  string `json:"email"`
	Avatar string `json:"avatar"`
}

type Data[T any] struct {
	Data T `json:"data"`
}

func Validate(req any) ([]string, bool) {
	v := validator.New()
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("fieldName"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	if err := v.Struct(req); err != nil {
		fieldError, oneofError, minError, defaultError := make([]string, 0), make([]string, 0), make([]string, 0), make([]string, 0)

		for _, err := range err.(validator.ValidationErrors) {
			switch err.Tag() {
			case "required":
				fieldError = append(fieldError, err.Field())
			case "oneof":
				oneofError = append(oneofError, err.Field())
			case "min", "len", "lte", "max":
				minError = append(minError, err.Field())
			default:
				defaultError = append(defaultError, err.Field())
			}
		}
		errs := []string{}
		if len(fieldError) > 0 {
			errs = append(errs, fmt.Sprintf("Los siguientes campos son requeridos: %v", strings.Join(fieldError, ", ")))
		}
		if len(oneofError) > 0 {
			errs = append(errs, fmt.Sprintf("Los siguientes campos no hacen match con el enumerado: %v", strings.Join(oneofError, ", ")))
		}
		if len(minError) > 0 {
			errs = append(errs, fmt.Sprintf("Los siguientes campos no cumplen con la longitud de caracteres requeridos: %v", strings.Join(minError, ", ")))
		}
		if len(defaultError) > 0 {
			errs = append(errs, fmt.Sprintf("Los siguientes campos son inválidos: %v", strings.Join(defaultError, ", ")))
		}
		return errs, false
	}
	return nil, true
}
