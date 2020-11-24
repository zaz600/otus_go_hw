package hw09_struct_validator //nolint:golint,stylecheck

import (
	"fmt"
	"log"
	"reflect"
	"strings"
)

const validateTag = "validate"

type ValidationError struct {
	Field string
	Err   error
}

func (e ValidationError) Error() string {
	return e.Err.Error()
}

func (e ValidationError) Unwrap() error {
	return e.Err
}

type ValidationErrors []ValidationError

type Validator interface {
	CanValidate(ruleName string) bool
	Validate(ruleValue string, fieldName string, fieldValue reflect.Value) *ValidationError
}

var validators map[string]Validator = map[string]Validator{
	"len":    LenValidator{},
	"min":    MinValidator{},
	"max":    MaxValidator{},
	"in":     InValidator{},
	"regexp": RegexpValidator{},
}

func findValidator(ruleName string) Validator {
	return validators[ruleName]
}

func (v ValidationErrors) Error() string {
	errorMessaged := make([]string, 0)
	for _, err := range v {
		errorMessaged = append(errorMessaged, fmt.Sprintf("field '%s': %s", err.Field, err.Err.Error()))
	}
	return strings.Join(errorMessaged, ", ")
}

func Validate(v interface{}) error {
	if v == nil {
		return nil
	}

	var errs ValidationErrors
	rValue := reflect.ValueOf(v)
	valueType := rValue.Type()

	if valueType.Kind() != reflect.Struct {
		log.Println("WARNING! Type must be struct")
		return nil
	}

	for i := 0; i < valueType.NumField(); i++ {
		field := valueType.Field(i)
		tag := field.Tag.Get(validateTag)
		if len(tag) == 0 {
			continue
		}

		fieldValue := rValue.Field(i)

		rules := strings.Split(tag, "|")
		for _, rule := range rules {
			ruleName, ruleValue, err := checkRule(rule)
			if err != nil {
				log.Printf("invalid tag format '%+v' for field '%s'", tag, field.Name)
				continue
			}

			validator := findValidator(ruleName)
			if validator != nil {
				err := validator.Validate(ruleValue, field.Name, fieldValue)
				if err != nil {
					errs = append(errs, *err)
				}
			} else {
				log.Printf("can't find validator for %s on field %s\n", ruleName, field.Name)
			}
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func checkRule(rule string) (string, string, error) {
	parts := strings.Split(rule, ":")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid tag format")
	}
	return parts[0], parts[1], nil
}

func ValidateSlice(validator Validator, ruleValue string, fieldName string, fieldValue reflect.Value) []ValidationError {
	errs := make([]ValidationError, 0)
	for i := 0; i < fieldValue.Len(); i++ {
		el := fieldValue.Index(i)
		if valErr := validator.Validate(ruleValue, fieldName, el); valErr != nil {
			errs = append(errs, *valErr)
		}
	}
	return errs
}
