package hw09_struct_validator //nolint:golint,stylecheck,dupl

import (
	"fmt"
	"reflect"
	"strconv"
)

type MaxValidator struct{}

func (v MaxValidator) CanValidate(ruleName string) bool {
	return ruleName == "max"
}

func (v MaxValidator) Validate(ruleValue string, fieldName string, fieldValue reflect.Value) *ValidationError {
	expected, err := strconv.Atoi(ruleValue)
	if err != nil {
		fmt.Printf("WARNING. Field: %s. invalid value='%s'. Must be int\n", fieldName, ruleValue)
		return nil
	}

	k := fieldValue.Kind()
	switch k {
	case reflect.Int:
		v := fieldValue.Int()
		if int(v) > expected {
			return &ValidationError{Field: fieldName, Err: fmt.Errorf("value must be less '%d'", expected)}
		}
		return nil

	case reflect.Slice:
		sliceErrs := ValidateSlice(v, ruleValue, fieldName, fieldValue)
		if len(sliceErrs) != 0 {
			return &ValidationError{Field: fieldName, Err: fmt.Errorf("all elements in slice must be less than '%d'", expected)}
		}

	default:
		fmt.Printf("WARNING. Field: %s. Invalid field type for len validator\n", fieldName)
		return nil
	}

	return nil
}
