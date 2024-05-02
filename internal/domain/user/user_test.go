package user

import (
	"errors"
	"testing"

	"syreclabs.com/go/faker"
)

var testLongPassword = faker.Internet().Password(PasswordMaxLength, 50)

func TestUser_ensureValidUsername(t *testing.T) {
	for _, test := range []struct {
		Name     string
		Input    string
		Expected error
	}{
		{Name: "with mixed letters and numbers", Input: "ragn4r0k", Expected: nil},
		{Name: "with only letters", Input: "ragnarok", Expected: nil},
		{Name: "with only numbers", Input: "23423", Expected: ErrUsernameStartWithLetter},
		{Name: "with empty username", Input: "", Expected: ErrEmptyUsername},
		{Name: "with a single letter", Input: "a", Expected: nil},
		{Name: "with a single number", Input: "4", Expected: ErrUsernameStartWithLetter},
		{Name: "with a non alphanumeric character", Input: "user-name", Expected: ErrUsernameOnlyAlphanumeric},
	} {
		t.Run(test.Name, func(t *testing.T) {
			got := ensureValidUsername(test.Input)
			if !errors.Is(got, test.Expected) {
				t.Errorf("got %v, expected %v", got, test.Expected)
			}
		})
	}
}

func TestUser_ensureValidPassword(t *testing.T) {
	for _, test := range []struct {
		Name     string
		Input    string
		Expected error
	}{
		{Name: "with valid password", Input: "password", Expected: nil},
		{Name: "with a short password", Input: "pass", Expected: ErrShortPassword},
		{Name: "with a long password", Input: testLongPassword, Expected: ErrLongPassword},
	} {
		t.Run(test.Name, func(t *testing.T) {
			got := ensureValidPassword(test.Input)
			if !errors.Is(got, test.Expected) {
				t.Errorf("got %v, expected %v", got, test.Expected)
			}
		})
	}
}
