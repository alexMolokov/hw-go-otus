package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
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
		in          interface{}
		expectedErr error
	}{
		{
			User{
				ID:    "123456789123456789123456789123456789",
				Name:  "Ok",
				Age:   18,
				Email: "alex@mail.ru",
				Role:  "admin",
				Phones: []string{
					"84957137708",
				},
				meta: json.RawMessage{1, 2, 3},
			},
			nil,
		},
		{
			App{
				Version: "1.3.6",
			},
			nil,
		},
		{
			Token{
				Header: []byte{1, 2, 3},
			},
			nil,
		},
		{
			Response{
				Code: 200,
			},
			nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)
			require.Equal(t, tt.expectedErr, err)
		})
	}

	testsFail := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			User{
				ID:    "123456789",
				Name:  "Fail",
				Age:   15,
				Email: "alexmail.ru",
				Role:  "admin",
				Phones: []string{
					"84957137708",
				},
			},
			ValidationErrors{
				{Field: "ID", Err: ErrValueIsInvalid},
				{Field: "Age", Err: ErrValueIsInvalid},
				{Field: "Email", Err: ErrValueIsInvalid},
			},
		},
	}

	for i, tt := range testsFail {
		t.Run(fmt.Sprintf("Fail case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)
			require.NotNil(t, err)
			errs, _ := err.(ValidationErrors) //nolint:errorlint
			require.Equal(t, 3, len(errs))
			expectedErrors, _ := tt.expectedErr.(ValidationErrors) //nolint:errorlint
			for i, e := range errs {
				require.True(t, errors.Is(e.Err, expectedErrors[i].Err))
			}
		})
	}
}
