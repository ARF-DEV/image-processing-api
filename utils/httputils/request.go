package httputils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func GetPageSize(r *http.Request, defaultPage, defaultSize int64) (int64, int64, error) {
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.ParseInt(pageStr, 10, 64)
	if errors.Is(err, strconv.ErrSyntax) {
		return 0, 0, ErrBadRequest
	}
	sizeStr := r.URL.Query().Get("size")
	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if errors.Is(err, strconv.ErrSyntax) {
		return 0, 0, ErrBadRequest
	}
	if page == 0 {
		page = defaultPage
	}
	if size == 0 {
		size = defaultSize
	}

	return page, size, nil
}

func GetURLParam[T any](r *http.Request, param string) (T, error) {
	var t T
	switch any(t).(type) {
	case int, int16, int32, int64:
		val := chi.URLParam(r, param)
		tVal, err := strconv.ParseInt(val, 10, 64)
		if errors.Is(err, strconv.ErrSyntax) {
			return t, ErrBadRequest
		}
		return any(tVal).(T), nil

	case string:
		val := chi.URLParam(r, param)
		if val == "" {
			return t, ErrBadRequest
		}
		return any(val).(T), nil
	}

	return t, fmt.Errorf("url param of type %T isn't implemented", t)
}

func ParseRequestBody(r *http.Request, dst any) error {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &dst); err != nil {
		return err
	}
	return nil
}
