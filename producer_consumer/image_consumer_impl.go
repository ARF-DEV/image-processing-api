package producerconsumer

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
	"log"
	"strings"
	"time"

	"github.com/ARF-DEV/image-processing-api/configs"
	"github.com/ARF-DEV/image-processing-api/model"
	"github.com/ARF-DEV/image-processing-api/repos/googlecloudstorage"
	"github.com/ARF-DEV/image-processing-api/repos/imagerepo"
	"github.com/disintegration/imaging"
	"github.com/rabbitmq/amqp091-go"
)

const (
	IMG_JPEG  string        = "jpeg"
	IMG_PNG   string        = "png"
	rateLimit time.Duration = time.Second / 20
)

type imageConvertFunc func(w io.Writer, image image.Image) error

type Consumer struct {
	ch        *amqp091.Channel
	conn      *amqp091.Connection
	imageRepo imagerepo.ImageRepo
	resource  googlecloudstorage.GoogleCloudStorageRepo
}

func NewConsumer(url string, imageRepo imagerepo.ImageRepo, resource googlecloudstorage.GoogleCloudStorageRepo) (*Consumer, error) {
	var err error
	consume := Consumer{
		imageRepo: imageRepo,
		resource:  resource,
	}
	consume.conn, err = amqp091.Dial(url)
	if err != nil {
		return nil, err
	}

	consume.ch, err = consume.conn.Channel()
	if err != nil {
		return nil, err
	}
	return &consume, nil
}

func (c *Consumer) RunConsumer(ctx context.Context, queueName string) {
	_, err := c.ch.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		log.Println("error when declaring queue: ", err)
		return
	}
	deliveryChan, err := c.ch.Consume(queueName, configs.GetConfig().QUEUE_NAME, false, false, false, false, nil)
	if err != nil {
		log.Println("error when consuming queue: ", err)
		return
	}

	ticker := time.NewTicker(rateLimit)
	defer ticker.Stop()

processMesssageLoop:
	for d := range deliveryChan {
		select {
		case <-ctx.Done():
			break processMesssageLoop
		case <-ticker.C:
			req := model.ImageTransformBrokerRequest{}
			if err := json.Unmarshal(d.Body, &req); err != nil {
				log.Println(err)
				d.Nack(false, false)
				continue
			}
			if err := c.TransformImage(ctx, req.ImageID, req.Req); err != nil {
				log.Println(err)
				d.Nack(false, false)
				continue
			}
			d.Ack(false)
		}
	}
}

func (c *Consumer) Close() {
	c.ch.Close()
	c.conn.Close()
}
func (s *Consumer) TransformImage(ctx context.Context, id int64, req model.ImageTransformRequestOpts) error {
	requestedImage, err := s.imageRepo.GetImage(ctx, id)
	if err != nil {
		return err
	}

	imageData, err := s.resource.LoadImage(ctx, requestedImage)
	if err != nil {
		return err
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
			return err
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
		return nil
	}
	decoder, ok := getDecodeFunctions()[imageData.Format]
	if !ok {
		return fmt.Errorf("decoder isn't implemented")
	}

	buf := bytes.Buffer{}
	if err := decoder(&buf, imageData.Image); err != nil {
		return err
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
		return err
	}

	newImage := model.Image{
		URL: url,
	}
	savedId, err := s.imageRepo.SaveImage(ctx, newImage)
	if err != nil {
		return err
	}

	newImage.ID = savedId
	// newImage := model.Image{}
	return nil
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
