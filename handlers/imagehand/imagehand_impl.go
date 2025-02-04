package imagehand

import (
	"net/http"

	"github.com/ARF-DEV/image-processing-api/model"
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

func (h *ImageHandlerImpl) GetImages(w http.ResponseWriter, r *http.Request) {
	page, limit, err := httputils.GetPageLimit(r, 1, 10)
	if err != nil {
		httputils.SendResponse(w, err.Error(), nil, nil, err)
		return
	}
	images, meta, err := h.imageServ.GetAllImage(r.Context(), page, limit)
	if err != nil {
		httputils.SendResponse(w, err.Error(), nil, nil, err)
		return
	}

	httputils.SendResponse(w, httputils.Success, images, meta, nil)
}

func (h *ImageHandlerImpl) GetImage(w http.ResponseWriter, r *http.Request) {
	imageId, err := httputils.GetURLParam[int64](r, "id")
	if err != nil {
		httputils.SendResponse(w, err.Error(), nil, nil, err)
		return
	}

	res, err := h.imageServ.GetImage(r.Context(), imageId)
	if err != nil {
		httputils.SendResponse(w, err.Error(), nil, nil, err)
		return
	}

	httputils.SendResponse(w, httputils.Success, res, nil, nil)
}

func (h *ImageHandlerImpl) TransformImage(w http.ResponseWriter, r *http.Request) {
	id, err := httputils.GetURLParam[int64](r, "id")
	if err != nil {
		httputils.SendResponse(w, err.Error(), nil, nil, err)
		return
	}
	transformReq := model.ImageTranformRequest{}
	if err := httputils.ParseRequestBody(r, &transformReq); err != nil {
		httputils.SendResponse(w, err.Error(), nil, nil, err)
		return
	}

	err = h.imageServ.TransformImageBroker(r.Context(), id, transformReq.Transform)
	if err != nil {
		httputils.SendResponse(w, err.Error(), nil, nil, err)
		return
	}
	httputils.SendResponse(w, httputils.Success, nil, nil, nil)
}
