package hw09_struct_validator //nolint:golint,stylecheck

import (
	"fmt"
	"reflect"
	"regexp"
)

type RegexpValidator struct{}

func (v RegexpValidator) CanValidate(ruleName string) bool {
	return ruleName == "regexp"
}

func (v RegexpValidator) Validate(ruleValue string, fieldName string, fieldValue reflect.Value) *ValidationError {
	re, err := regexp.Compile(ruleValue)
	if err != nil {
		fmt.Printf("WARNING. Field: %s. invalid regexp value='%s'.\n", fieldName, ruleValue)
		return nil
	}

	k := fieldValue.Kind()
	switch k {
	case reflect.String:
		s := fieldValue.String()
		if s == "" {
			return nil
		}
		if !re.MatchString(s) {
			return &ValidationError{Field: fieldName, Err: fmt.Errorf("must match regexp '%s'", ruleValue)}
		}

	case reflect.Slice:
		sliceErrs := ValidateSlice(v, ruleValue, fieldName, fieldValue)
		if len(sliceErrs) != 0 {
			return &ValidationError{Field: fieldName, Err: fmt.Errorf("all elements in slice must match regexp %s", ruleValue)}
		}

	default:
		fmt.Printf("WARNING. Field: %s. Invalid field type for 'regexp' validator\n", fieldName)
		return nil
	}

	return nil
}
