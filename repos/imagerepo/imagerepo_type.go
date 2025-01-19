package imagerepo

import (
	"context"

	"github.com/ARF-DEV/image-processing-api/model"
)

type ImageRepo interface {
	SaveImage(ctx context.Context, image model.Image) error
}
