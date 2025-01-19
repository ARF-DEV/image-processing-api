package model

import "io"

type UploadImageRequest struct {
	Reader io.Reader
	Name   string
}

type Image struct {
	ID  int64  `db:"id"`
	URL string `db:"url"`
}
