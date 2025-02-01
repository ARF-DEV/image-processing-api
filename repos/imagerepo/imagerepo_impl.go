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

func (r ImageRepoImpl) SaveImage(ctx context.Context, image model.Image) (int64, error) {
	sq := squirrel.Insert("images").Columns("url").Values(image.URL).Suffix("RETURNING id")
	query, args, err := sq.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return 0, err
	}

	stmt, err := r.db.PreparexContext(ctx, query)
	if err != nil {
		return 0, err
	}

	var id int64
	err = stmt.QueryRowxContext(ctx, args...).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (r *ImageRepoImpl) GetImages(ctx context.Context, page, limit int64) ([]model.Image, error) {
	offset := (page - 1) * limit
	sq := squirrel.Select("id", "url").From("images").Limit(uint64(limit)).Offset(uint64(offset))
	query, args, err := sq.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	stmt, err := r.db.PreparexContext(ctx, query)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.QueryxContext(ctx, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var images []model.Image

	for rows.Next() {
		var image model.Image
		if err := rows.StructScan(&image); err != nil {
			return nil, err
		}
		images = append(images, image)
	}

	return images, nil
}

func (r *ImageRepoImpl) CountImages(ctx context.Context) (int64, error) {
	sq := squirrel.Select("count(id)").From("images")

	query, args, err := sq.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return 0, err
	}

	stmt, err := r.db.PreparexContext(ctx, query)
	if err != nil {
		return 0, err
	}

	var count int64
	if err := stmt.QueryRowxContext(ctx, args...).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (r *ImageRepoImpl) GetImage(ctx context.Context, id int64) (model.Image, error) {
	sq := squirrel.Select("id", "url").From("images").Where(squirrel.Eq{"id": id})

	query, args, err := sq.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return model.Image{}, err
	}

	stmt, err := r.db.PreparexContext(ctx, query)
	if err != nil {
		return model.Image{}, err
	}

	var image model.Image
	if err := stmt.QueryRowxContext(ctx, args...).StructScan(&image); err != nil {
		return model.Image{}, err
	}

	return image, nil
}
