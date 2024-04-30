package user

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID
	Username string
	Password string
}

func New(id uuid.UUID, username, password string) User {
	return User{ID: id, Username: username, Password: password}
}

func Create(username, password string) User {
	return New(uuid.New(), username, password)
}
