package model

import (
	"fmt"
	"image"
	"io"
	"strings"

	"github.com/ARF-DEV/image-processing-api/configs"
)

type UploadImageRequest struct {
	Reader io.Reader
	Name   string
}

type Image struct {
	ID  int64  `db:"id"`
	URL string `db:"url"`
}

func (i Image) ToImageResponse(cfg *configs.Config) ImageResponse {
	image := ImageResponse{
		ID:  i.ID,
		URL: i.URL,
	}

	if image.URL != "" {
		image.URL = fmt.Sprintf("%s%s", cfg.GOOGLE_STORAGE_URL, image.URL)
	}
	return image
}

type ImageResponse struct {
	ID  int64  `json:"id"`
	URL string `json:"url"`
}

type ImageResponses []ImageResponse

type Images []Image

func (i Images) ToImageResponses(cfg *configs.Config) ImageResponses {
	var imageResponses ImageResponses
	for _, image := range i {
		imageRes := image.ToImageResponse(cfg)
		imageResponses = append(imageResponses, imageRes)
	}
	return imageResponses

}

type ImageTransformRequestOpts struct {
	ResizeTransform ResizeTransformRequest `json:"resize"`
	CropTransform   CropTransformRequest   `json:"crop"`
	Rotate          float64                `json:"rotate"`
	Format          string                 `json:"format"`
	Filters         FilterTransformRequest `json:"filters"`
}

type ImageTranformRequest struct {
	Transform ImageTransformRequestOpts `json:"transformations"`
}

func (i *ImageTransformRequestOpts) GenerateStr() string {

	s := []string{}
	if i.ResizeTransform != (ResizeTransformRequest{}) {
		s = append(s, "resized")
	}
	if i.CropTransform != (CropTransformRequest{}) {
		s = append(s, "cropped")
	}
	if i.Rotate > 0 {
		s = append(s, "rotated")
	}
	if i.Format != "" {
		s = append(s, "formated")
	}
	if i.Filters != (FilterTransformRequest{}) {
		s = append(s, "filtered")
	}

	return strings.Join(s, "-")
}

type ResizeTransformRequest struct {
	Width  int64 `json:"width"`
	Height int64 `json:"height"`
}
type CropTransformRequest struct {
	X      int64 `json:"x"`
	Y      int64 `json:"y"`
	Width  int64 `json:"width"`
	Height int64 `json:"height"`
}
type FilterTransformRequest struct {
	Grayscale bool `json:"grayscale"`
	Sepia     bool `json:"sepia"`
}

func (i *Image) GetBucket() string {
	strSplits := strings.Split(i.URL, "/")
	return strSplits[1]
}

func (i *Image) GetObject() string {
	strSplits := strings.Split(i.URL, "/")
	return strSplits[2]
}

type ImageInfo struct {
	Image  image.Image
	Format string
}

type ImageTransformBrokerRequest struct {
	Req     ImageTransformRequestOpts `json:"opts"`
	ImageID int64                     `json:"image_id"`
}
