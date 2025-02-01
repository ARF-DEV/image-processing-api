package imageserv

import (
	"context"
	"mime/multipart"

	"github.com/ARF-DEV/image-processing-api/model"
)

type ImageServ interface {
	UploadImage(ctx context.Context, file multipart.File, header *multipart.FileHeader) error
	GetAllImage(ctx context.Context, page int64, limit int64) (model.ImageResponses, *model.Meta, error)
	GetImage(ctx context.Context, id int64) (model.ImageResponse, error)
	TransformImage(ctx context.Context, id int64, req model.ImageTransformRequestOpts) (model.ImageResponse, error)
}
