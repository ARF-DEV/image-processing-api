package handlers

import (
	"net/http"

	"github.com/ARF-DEV/image-processing-api/handlers/imagehand"
	"github.com/ARF-DEV/image-processing-api/handlers/userhand"
	"github.com/go-chi/chi/v5"
)

func CreateHandlers(user userhand.UserHandler, image imagehand.ImageHandler) http.Handler {
	r := chi.NewRouter()

	r.Post("/register", user.Register)
	r.Post("/login", user.Login)

	r.Route("/images", func(r chi.Router) {
		// r.Use(middleware.Authenticate)
		r.Get("/", image.GetImages)
		r.Post("/", image.UploadImage)

		r.Get("/{id}", image.GetImage)
		r.Post("/{id}/transform", image.TransformImage)
	})

	return r
}
