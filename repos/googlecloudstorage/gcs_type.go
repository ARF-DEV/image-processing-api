package googlecloudstorage

import (
	"context"

	"github.com/ARF-DEV/image-processing-api/model"
)

type GoogleCloudStorageRepo interface {
	CreateBucket(ctx context.Context) error
	UploadImage(ctx context.Context, req model.UploadImageRequest) (string, error)
	Close()
}
