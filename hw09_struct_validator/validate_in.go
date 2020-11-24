package hw09_struct_validator //nolint:golint,stylecheck

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type InValidator struct{}

func (v InValidator) CanValidate(ruleName string) bool {
	return ruleName == "in"
}

func (v InValidator) Validate(ruleValue string, fieldName string, fieldValue reflect.Value) *ValidationError {
	expectedValues := strings.Split(ruleValue, ",")

	k := fieldValue.Kind()
	switch k {
	case reflect.Int:
		expectedInts, err := sliceStringToInt(expectedValues)
		if err != nil {
			fmt.Printf("WARNING. Field: %s. validation for int field must be (slice) int\n", fieldName)
			return nil
		}

		actual := int(fieldValue.Int())
		for _, expected := range expectedInts {
			if actual == expected {
				return nil
			}
		}
		return &ValidationError{Field: fieldName, Err: fmt.Errorf("value must be in '%+v'", expectedInts)}

	case reflect.String:
		actual := fieldValue.String()
		for _, expected := range expectedValues {
			if actual == expected {
				return nil
			}
		}
		return &ValidationError{Field: fieldName, Err: fmt.Errorf("value must be in '%+v'", expectedValues)}

	case reflect.Slice:
		sliceErrs := ValidateSlice(v, ruleValue, fieldName, fieldValue)
		if len(sliceErrs) != 0 {
			return &ValidationError{Field: fieldName, Err: fmt.Errorf("all elements in slice must be in '%+v'", expectedValues)}
		}

	default:
		fmt.Printf("WARNING. Field: %s. Invalid field type for 'in' validator\n", fieldName)
		return nil
	}

	return nil
}

func sliceStringToInt(slice []string) ([]int, error) {
	result := make([]int, 0)
	for _, s := range slice {
		value, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("can't convert %s to int", s)
		}
		result = append(result, value)
	}
	return result, nil
}
