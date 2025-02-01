package imagerepo

import (
	"context"

	"github.com/ARF-DEV/image-processing-api/model"
)

type ImageRepo interface {
	SaveImage(ctx context.Context, image model.Image) (int64, error)
	GetImages(ctx context.Context, page int64, limit int64) ([]model.Image, error)
	CountImages(ctx context.Context) (int64, error)
	GetImage(ctx context.Context, id int64) (model.Image, error)
}
