package hw09_struct_validator //nolint:golint,stylecheck

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"syreclabs.com/go/faker"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Foo struct {
		Bar []string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			name: "no tags",
			in: Token{
				Header:    []byte{10, 20, 100},
				Payload:   []byte{50, 60, 70},
				Signature: []byte{2, 3, 4},
			},
			expectedErr: nil,
		},

		{
			name: "valid len string",
			in:   App{Version: faker.Number().Number(5)},
		},

		{
			name: "invalid len string (greater)",
			in:   App{Version: faker.Number().Number(6)},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Version", Err: errors.New("length must be 5")},
			},
		},

		{
			name: "invalid len string (lower)",
			in:   App{Version: faker.Number().Number(4)},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Version", Err: errors.New("length must be 5")},
			},
		},

		{
			name: "valid len slice",
			in:   Foo{Bar: []string{"12345", "12345"}},
		},
		{
			name: "invalid len slice",
			in:   Foo{Bar: []string{"123456", "123456"}},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Bar", Err: errors.New("all elements in slice must have length 5")},
			},
		},

		{
			name: "valid in int",
			in:   Response{Code: 200, Body: "message"},
		},

		{
			name: "all validators valid",
			in: User{
				ID:    faker.Number().Number(36),
				Name:  "foo",
				Age:   30,
				Email: "foo@example.com",
				Role:  "admin",
				Phones: []string{
					faker.Number().Number(11),
					faker.Number().Number(11),
				},
			},
		},

		{
			name: "all validators invalid",
			in: User{
				ID:    faker.Number().Number(3),
				Name:  "foo",
				Age:   60,
				Email: "foo example.com",
				Role:  "superadmin",
				Phones: []string{
					faker.Number().Number(9),
					faker.Number().Number(8),
				},
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "ID", Err: errors.New("length must be 36")},
				ValidationError{Field: "Age", Err: errors.New("value must be less '50'")},
				ValidationError{Field: "Email", Err: errors.New("must match regexp '^\\w+@\\w+\\.\\w+$'")},
				ValidationError{Field: "Role", Err: errors.New("value must be in '[admin stuff]'")},
				ValidationError{Field: "Phones", Err: errors.New("all elements in slice must have length 11")},
			},
		},

		{
			name: "nil",
			in:   nil,
		},

		{
			name: "not a struct",
			in:   "foo",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d: %s", i, tt.name), func(t *testing.T) {
			err := Validate(tt.in)

			if tt.expectedErr != nil {
				require.Equal(t, tt.expectedErr, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
