package hw09_struct_validator //nolint:golint,stylecheck

import (
	"fmt"
	"reflect"
	"strconv"
)

type LenValidator struct{}

func (v LenValidator) CanValidate(ruleName string) bool {
	return ruleName == "len"
}

func (v LenValidator) Validate(ruleValue string, fieldName string, fieldValue reflect.Value) *ValidationError {
	expectedLen, err := strconv.Atoi(ruleValue)
	if err != nil {
		fmt.Printf("WARNING. Field: %s. invalid len value='%s'. Must be int\n", fieldName, ruleValue)
		return nil
	}

	k := fieldValue.Kind()
	switch k {
	case reflect.String:
		s := fieldValue.String()
		if len(s) != expectedLen {
			return &ValidationError{Field: fieldName, Err: fmt.Errorf("length must be %v", expectedLen)}
		}

	case reflect.Slice:
		sliceErrs := ValidateSlice(v, ruleValue, fieldName, fieldValue)
		if len(sliceErrs) != 0 {
			return &ValidationError{Field: fieldName, Err: fmt.Errorf("all elements in slice must have length %d", expectedLen)}
		}

	default:
		fmt.Printf("WARNING. Field: %s. Invalid field type for len validator\n", fieldName)
		return nil
	}

	return nil
}
