package imageserv

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"mime/multipart"
	"strings"

	"github.com/ARF-DEV/image-processing-api/configs"
	"github.com/ARF-DEV/image-processing-api/model"
	producerconsumer "github.com/ARF-DEV/image-processing-api/producer_consumer"
	"github.com/ARF-DEV/image-processing-api/repos/googlecloudstorage"
	"github.com/ARF-DEV/image-processing-api/repos/imagerepo"
	"github.com/disintegration/imaging"
)

const (
	IMG_JPEG string = "jpeg"
	IMG_PNG  string = "png"
)

type imageConvertFunc func(w io.Writer, image image.Image) error

type ImageServImpl struct {
	resource  googlecloudstorage.GoogleCloudStorageRepo
	imageRepo imagerepo.ImageRepo
	producer  *producerconsumer.Producer
}

func New(resource googlecloudstorage.GoogleCloudStorageRepo, imageRepo imagerepo.ImageRepo, producer *producerconsumer.Producer) ImageServ {
	return &ImageServImpl{
		resource:  resource,
		imageRepo: imageRepo,
		producer:  producer,
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

	if _, err := s.imageRepo.SaveImage(ctx, model.Image{
		URL: url,
	}); err != nil {
		return err
	}

	return nil
}

func (s *ImageServImpl) GetAllImage(ctx context.Context, page int64, limit int64) (model.ImageResponses, *model.Meta, error) {
	images, err := s.imageRepo.GetImages(ctx, page, limit)
	if err != nil {
		return nil, nil, err
	}

	total, err := s.imageRepo.CountImages(ctx)
	if err != nil {
		return nil, nil, err
	}

	meta := model.Meta{
		Page:      page,
		Limit:     limit,
		TotalData: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(limit))),
	}
	return model.Images(images).ToImageResponses(configs.GetConfig()), &meta, nil
}

func (s *ImageServImpl) GetImage(ctx context.Context, id int64) (model.ImageResponse, error) {
	image, err := s.imageRepo.GetImage(ctx, id)
	if err != nil {
		return model.ImageResponse{}, err
	}
	return image.ToImageResponse(configs.GetConfig()), nil
}

func (s *ImageServImpl) TransformImage(ctx context.Context, id int64, req model.ImageTransformRequestOpts) (model.ImageResponse, error) {
	requestedImage, err := s.imageRepo.GetImage(ctx, id)
	if err != nil {
		return model.ImageResponse{}, err
	}

	imageData, err := s.resource.LoadImage(ctx, requestedImage)
	if err != nil {
		return model.ImageResponse{}, err
	}
	transformed := false
	if req.CropTransform != (model.CropTransformRequest{}) {
		transformed = true
		imageData.Image = CropImage(imageData.Image, req.CropTransform)
	}

	if req.Format != "" {
		transformed = true
		imageData.Image, err = ChangeImageFormat(imageData.Image, req.Format)
		if err != nil {
			return model.ImageResponse{}, err
		}
	}
	if req.Filters != (model.FilterTransformRequest{}) {
		transformed = true
		if req.Filters.Grayscale {
			imageData.Image = GrayscaleFilterImage(imageData.Image)
		}
		if req.Filters.Sepia {
			imageData.Image = SepiaFilterImage(imageData.Image)
		}
	}

	if req.ResizeTransform != (model.ResizeTransformRequest{}) {
		transformed = true
		imageData.Image = ResizeImage(imageData.Image, req.ResizeTransform)
	}
	if req.Rotate > 0 {
		transformed = true
		imageData.Image = RotateImage(imageData.Image, req.Rotate)
	}

	if !transformed {
		return requestedImage.ToImageResponse(configs.GetConfig()), nil
	}
	decoder, ok := getDecodeFunctions()[imageData.Format]
	if !ok {
		return model.ImageResponse{}, fmt.Errorf("decoder isn't implemented")
	}

	buf := bytes.Buffer{}
	if err := decoder(&buf, imageData.Image); err != nil {
		return model.ImageResponse{}, err
	}

	uploadReq := model.UploadImageRequest{
		Reader: &buf,
	}
	strSplit := strings.Split(requestedImage.GetObject(), ".")
	fileName := strSplit[0]
	fileExtentions := strSplit[1]
	uploadReq.Name = fmt.Sprintf("%s:%s.%s", fileName, req.GenerateStr(), fileExtentions)
	url, err := s.resource.UploadImage(ctx, uploadReq)
	if err != nil {
		return model.ImageResponse{}, err
	}

	newImage := model.Image{
		URL: url,
	}
	savedId, err := s.imageRepo.SaveImage(ctx, newImage)
	if err != nil {
		return model.ImageResponse{}, err
	}

	newImage.ID = savedId
	// newImage := model.Image{}
	return newImage.ToImageResponse(configs.GetConfig()), nil
}

func CropImage(imageData image.Image, cropReq model.CropTransformRequest) image.Image {
	newImage := image.NewRGBA(imageData.Bounds())
	draw.Draw(newImage, newImage.Bounds(), imageData, newImage.Rect.Min, draw.Src)

	return newImage.SubImage(image.Rect(int(cropReq.X), int(cropReq.Y), int(cropReq.Width+cropReq.X), int(cropReq.Height+cropReq.Y)))
}

func RotateImage(imageData image.Image, rotateReq float64) image.Image {
	return imaging.Rotate(imageData, rotateReq, color.Black)
}

func ResizeImage(imageData image.Image, resizeReq model.ResizeTransformRequest) image.Image {
	// resize using nearest neighbour algorithm
	// ref: https://medium.com/@chathuragunasekera/image-resampling-algorithms-for-pixel-manipulation-bee65dda1488
	heightScale := float64(imageData.Bounds().Dy()) / float64(resizeReq.Height)
	widthScale := float64(imageData.Bounds().Dx()) / float64(resizeReq.Width)

	newImage := image.NewRGBA(image.Rect(0, 0, int(resizeReq.Width), int(resizeReq.Height)))
	for y := 0; y < newImage.Bounds().Dy(); y++ {
		for x := 0; x < newImage.Bounds().Dx(); x++ {
			xCoords := x * int(widthScale)
			yCoords := y * int(heightScale)

			newImage.Set(x, y, imageData.At(xCoords, yCoords))
		}
	}

	return newImage
}
func ChangeImageFormat(imageData image.Image, targetFormat string) (image.Image, error) {
	decoder, found := getDecodeFunctions()[targetFormat]
	if !found {
		return nil, fmt.Errorf("image decoder for %s not found", targetFormat)
	}
	var buf bytes.Buffer
	if err := decoder(&buf, imageData); err != nil {
		return nil, err
	}

	newFormatImage, _, err := image.Decode(&buf)
	if err != nil {
		return nil, err
	}

	return newFormatImage, nil
}
func GrayscaleFilterImage(imageData image.Image) image.Image {
	newImage := image.NewRGBA(imageData.Bounds())
	for y := 0; y < imageData.Bounds().Dy(); y++ {
		for x := 0; x < imageData.Bounds().Dx(); x++ {
			r, g, b, a := imageData.At(x, y).RGBA()
			// ref: https://www.johndcook.com/blog/2009/08/24/algorithms-convert-color-grayscale/
			grayVal := 0.21*float64(r) + 0.72*float64(g) + 0.07*float64(b)
			newImage.SetRGBA64(x, y, color.RGBA64{uint16(grayVal), uint16(grayVal), uint16(grayVal), uint16(a)})
		}
	}
	return newImage
}

func SepiaFilterImage(imageData image.Image) image.Image {
	newImage := image.NewRGBA(imageData.Bounds())
	for y := 0; y < imageData.Bounds().Dy(); y++ {
		for x := 0; x < imageData.Bounds().Dx(); x++ {
			r, g, b, a := imageData.At(x, y).RGBA()
			tr := 0.393*float64(r) + 0.769*float64(g) + 0.189*float64(b)
			tg := 0.349*float64(r) + 0.686*float64(g) + 0.168*float64(b)
			tb := 0.272*float64(r) + 0.534*float64(g) + 0.131*float64(b)
			newRGB := color.RGBA64{A: uint16(a)}
			rMax, gMax, bMax, _ := color.White.RGBA()
			if tr > float64(rMax) {
				newRGB.R = uint16(r)
			} else {
				newRGB.R = uint16(tr)
			}
			if tg > float64(gMax) {
				newRGB.G = uint16(g)
			} else {
				newRGB.G = uint16(tg)
			}
			if tb > float64(bMax) {
				newRGB.B = uint16(b)
			} else {
				newRGB.B = uint16(tb)
			}

			newImage.SetRGBA64(x, y, newRGB)
		}
	}
	return newImage
}
func getDecodeFunctions() map[string]imageConvertFunc {
	return map[string]imageConvertFunc{
		IMG_JPEG: func(w io.Writer, image image.Image) error {
			return jpeg.Encode(w, image, nil)
		},
		IMG_PNG: png.Encode,
	}
}

func (s *ImageServImpl) TransformImageBroker(ctx context.Context, id int64, req model.ImageTransformRequestOpts) error {
	data, err := json.Marshal(model.ImageTransformBrokerRequest{
		ImageID: id,
		Req:     req,
	})
	if err != nil {
		return err
	}
	err = s.producer.PublishCtx(ctx, configs.GetConfig().QUEUE_NAME, data)
	if err != nil {
		return err
	}
	return nil
}
