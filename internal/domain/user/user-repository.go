package user

import "context"

type Repository interface {
	Save(context.Context, User) (User, error)
}
