package imagehand

import (
	"net/http"

	"github.com/ARF-DEV/image-processing-api/services/imageserv"
	"github.com/ARF-DEV/image-processing-api/utils/httputils"
)

type ImageHandlerImpl struct {
	imageServ imageserv.ImageServ
}

func New(imageServ imageserv.ImageServ) ImageHandler {
	return &ImageHandlerImpl{
		imageServ: imageServ,
	}
}

func (h *ImageHandlerImpl) UploadImage(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(1024); err != nil {
		httputils.SendResponse(w, err.Error(), nil, nil, err)
		return
	}

	img, header, err := r.FormFile("image")
	if err != nil {
		httputils.SendResponse(w, err.Error(), nil, nil, err)
		return
	}

	if err := h.imageServ.UploadImage(r.Context(), img, header); err != nil {
		httputils.SendResponse(w, err.Error(), nil, nil, err)
		return
	}

	httputils.SendResponse(w, httputils.Success, nil, nil, nil)
}
