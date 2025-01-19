package imagerepo

import (
	"context"

	"github.com/ARF-DEV/image-processing-api/model"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type ImageRepoImpl struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) ImageRepo {
	return &ImageRepoImpl{db: db}
}

func (r ImageRepoImpl) SaveImage(ctx context.Context, image model.Image) error {
	sq := squirrel.Insert("images").Columns("url").Values(image.URL)
	query, args, err := sq.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}

	stmt, err := r.db.PreparexContext(ctx, query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, args...)
	if err != nil {
		return err
	}

	return nil
}
