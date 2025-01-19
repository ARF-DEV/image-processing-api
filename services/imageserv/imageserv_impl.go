package imageserv

import (
	"context"
	"mime/multipart"

	"github.com/ARF-DEV/image-processing-api/model"
	"github.com/ARF-DEV/image-processing-api/repos/googlecloudstorage"
	"github.com/ARF-DEV/image-processing-api/repos/imagerepo"
)

type ImageServImpl struct {
	resource  googlecloudstorage.GoogleCloudStorageRepo
	imageRepo imagerepo.ImageRepo
}

func New(resource googlecloudstorage.GoogleCloudStorageRepo, imageRepo imagerepo.ImageRepo) ImageServ {
	return &ImageServImpl{
		resource:  resource,
		imageRepo: imageRepo,
	}
}

func (s *ImageServImpl) UploadImage(ctx context.Context, file multipart.File, header *multipart.FileHeader) error {
	url, err := s.resource.UploadImage(ctx, model.UploadImageRequest{
		Name:   header.Filename,
		Reader: file,
	})
	if err != nil {
		return err
	}

	if err := s.imageRepo.SaveImage(ctx, model.Image{
		URL: url,
	}); err != nil {
		return err
	}

	return nil
}
