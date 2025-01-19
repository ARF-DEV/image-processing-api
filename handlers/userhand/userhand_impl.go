package userhand

import (
	"net/http"

	"github.com/ARF-DEV/image-processing-api/model"
	"github.com/ARF-DEV/image-processing-api/services/userserv"
	"github.com/ARF-DEV/image-processing-api/utils/httputils"
)

type UserHandlerImpl struct {
	userServ userserv.UserServ
}

func New(userServ userserv.UserServ) UserHandler {
	return &UserHandlerImpl{userServ: userServ}
}

func (h *UserHandlerImpl) Login(w http.ResponseWriter, r *http.Request) {
	req := model.LoginRegisterRequest{}
	if err := httputils.ParseRequestBody(r, &req); err != nil {
		httputils.SendResponse(w, err.Error(), nil, nil, err)
		return
	}
	res, err := h.userServ.Login(r.Context(), model.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		httputils.SendResponse(w, err.Error(), nil, nil, err)
		return
	}

	httputils.SendResponse(w, httputils.Success, res, nil, nil)
}

func (h *UserHandlerImpl) Register(w http.ResponseWriter, r *http.Request) {
	req := model.LoginRegisterRequest{}
	if err := httputils.ParseRequestBody(r, &req); err != nil {
		httputils.SendResponse(w, err.Error(), nil, nil, err)
		return
	}
	err := h.userServ.Register(r.Context(), model.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		httputils.SendResponse(w, err.Error(), nil, nil, err)
		return
	}

	httputils.SendResponse(w, httputils.Success, nil, nil, nil)
}
