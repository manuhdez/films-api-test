package mocks

import "github.com/stretchr/testify/mock"

type MockPasswordHasher struct {
	mock.Mock
}

func (m *MockPasswordHasher) Hash(password string) (string, error) {
	args := m.Called(password)
	return args.Get(0).(string), args.Error(1)
}

func (m *MockPasswordHasher) Compare(hash, password string) bool {
	args := m.Called(hash, password)
	return args.Bool(0)
}
