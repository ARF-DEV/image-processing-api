package userserv

import (
	"context"

	"github.com/ARF-DEV/image-processing-api/model"
)

type UserServ interface {
	Login(ctx context.Context, user model.User) (model.AutheticationResponse, error)
	Register(ctx context.Context, user model.User) error
}
