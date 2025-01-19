package userrepo

import (
	"context"

	"github.com/ARF-DEV/image-processing-api/model"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type UserRepoImpl struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) UserRepo {
	return &UserRepoImpl{db: db}
}

func (r *UserRepoImpl) CreateUser(ctx context.Context, user model.User) error {
	sq := squirrel.Insert("users").Columns("email", "password").Values(user.Email, user.Password)
	query, args, err := sq.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}

	stmt, err := r.db.PreparexContext(ctx, query)
	if err != nil {
		return err
	}

	if _, err := stmt.ExecContext(ctx, args...); err != nil {
		return err
	}

	return nil
}

func (r *UserRepoImpl) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	sq := squirrel.Select("id", "email", "password").From("users").Limit(1)
	query, args, err := sq.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return model.User{}, err
	}

	stmt, err := r.db.PreparexContext(ctx, query)
	if err != nil {
		return model.User{}, err
	}

	var res model.User
	if err := stmt.QueryRowxContext(ctx, args...).StructScan(&res); err != nil {
		return model.User{}, err
	}
	return res, nil
}
