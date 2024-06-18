package mocks

import "github.com/stretchr/testify/mock"

type MockTokenGenerator struct {
	mock.Mock
}

func (m *MockTokenGenerator) Generate(id string) (string, error) {
	args := m.Called(id)
	return args.Get(0).(string), args.Error(1)
}

func (m *MockTokenGenerator) Validate(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
