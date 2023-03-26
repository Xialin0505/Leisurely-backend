package utils

import (
	"fmt"

	_ "github.com/go-playground/validator"
	"github.com/go-playground/validator/v10"
)

func constructFormError(err validator.FieldError, errFields *[]string, errMsgs *[]string) {
	errMsg := err.StructField()
	*errFields = append(*errFields, errMsg)
	if err.Tag() == "required" {
		errMsg += " is required"
		*errMsgs = append(*errMsgs, errMsg+", but got: "+fmt.Sprintf("%v", err.Value()))
		return
	} else {
		errMsg += " should be " + err.Tag()
	}
	if err.Param() != "" {
		errMsg += " " + err.Param()
	}
	*errMsgs = append(*errMsgs, errMsg+", but got: "+fmt.Sprintf("%v", err.Value()))
}

func ValidateDTO(form interface{}) map[string]string {
	var errFields []string
	var errMsgs []string

	validate := validator.New()
	if err := validate.Struct(form); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			constructFormError(err, &errFields, &errMsgs)
		}
	}

	if len(errFields) == 0 || len(errMsgs) == 0 {
		return nil
	}

	res := make(map[string]string)

	for i, field := range errFields {
		res[field] = errMsgs[i]
	}

	return res
}
