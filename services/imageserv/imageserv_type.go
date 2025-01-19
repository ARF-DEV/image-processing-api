package imageserv

import (
	"context"
	"mime/multipart"
)

type ImageServ interface {
	UploadImage(ctx context.Context, file multipart.File, header *multipart.FileHeader) error
}
