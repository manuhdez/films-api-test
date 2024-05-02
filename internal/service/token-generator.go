package service

type TokenGenerator interface {
	Generate(id string) (string, error)
	Validate(token string) error
}
