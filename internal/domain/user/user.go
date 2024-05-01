package user

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/google/uuid"
)

const (
	PasswordMinLength = 8
	PasswordMaxLength = 24
)

var (
	ErrEmptyUsername            = errors.New("username must not be empty")
	ErrUsernameOnlyAlphanumeric = errors.New("username must only contain alphanumeric characters")
	ErrUsernameStartWithLetter  = errors.New("username must start with a letter")
	ErrShortPassword            = fmt.Errorf("password must be at least %d characters", PasswordMinLength)
	ErrLongPassword             = fmt.Errorf("password must not be longer than %d characters", PasswordMaxLength)
)

type User struct {
	ID       uuid.UUID
	Username string
	Password string
}

func New(id uuid.UUID, username, password string) (User, error) {
	if err := ensureValidUsername(username); err != nil {
		return User{}, err
	}

	if err := ensureValidPassword(password); err != nil {
		return User{}, err
	}

	return User{
		ID:       id,
		Username: username,
		Password: password,
	}, nil
}

func Create(username, password string) (User, error) {
	return New(uuid.New(), username, password)
}

func ensureValidUsername(value string) error {
	// Should have at least one character
	if len(value) == 0 {
		return ErrEmptyUsername
	}

	// Should only include alphanumeric characters
	if !regexp.MustCompile("^[a-zA-Z0-9]+$").MatchString(value) {
		return ErrUsernameOnlyAlphanumeric
	}

	// Should start with a letter
	if !regexp.MustCompile("^[a-zA-Z]").MatchString(value) {
		return ErrUsernameStartWithLetter
	}

	return nil
}

func ensureValidPassword(value string) error {
	if len(value) < PasswordMinLength {
		return ErrShortPassword
	}

	if len(value) > PasswordMaxLength {
		return ErrLongPassword
	}

	return nil
}
