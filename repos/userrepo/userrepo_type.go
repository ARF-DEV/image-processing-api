package userrepo

import (
	"context"

	"github.com/ARF-DEV/image-processing-api/model"
)

type UserRepo interface {
	CreateUser(ctx context.Context, user model.User) error
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
}
