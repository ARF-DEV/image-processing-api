package imagehand

import "net/http"

type ImageHandler interface {
	UploadImage(w http.ResponseWriter, r *http.Request)
	GetImages(w http.ResponseWriter, r *http.Request)
	GetImage(w http.ResponseWriter, r *http.Request)
	TransformImage(w http.ResponseWriter, r *http.Request)
}
