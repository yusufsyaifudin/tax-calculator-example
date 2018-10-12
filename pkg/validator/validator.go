package validator

import (
	"fmt"
	"strings"

	"github.com/serenize/snaker"
	libValidator "gopkg.in/go-playground/validator.v9"
)

var v *libValidator.Validate

func init() {
	v = libValidator.New()
}

// Errors wraps validation errors
// It maps a field and the errors for the field as string
type Errors struct {
	Data map[string][]string
}

func (e *Errors) String() string {
	var errorMessage []string
	for key, msg := range e.Data {
		errorMessage = append(errorMessage, fmt.Sprintf("%s: %s", key, strings.Join(msg, ",")))
	}

	return strings.Join(errorMessage, "; ")
}

// Validate validates a struct
func Validate(f interface{}) *Errors {
	err := v.Struct(f)
	if err == nil {
		return nil
	}

	var errorData = make(map[string][]string)
	for _, err := range err.(libValidator.ValidationErrors) {
		field := snaker.CamelToSnake(err.Field())
		msg := err.ActualTag()
		if _, ok := errorData[field]; !ok {
			errorData[field] = []string{}
		}

		errorData[field] = append(errorData[field], msg)
	}
	return &Errors{
		Data: errorData,
	}
}

// AddError adds error pair (field, msg) to current Errors instance
func AddError(errs *Errors, field, msg string) {
	if errs == nil {
		return
	}

	var errorData = errs.Data
	if errorData == nil {
		errorData = make(map[string][]string)
	}

	if _, exist := errorData[field]; !exist {
		errorData[field] = []string{}
	}

	errorData[field] = append(errorData[field], msg)

	// set new value
	errs.Data = errorData
}
