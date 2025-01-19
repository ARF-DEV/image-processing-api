package httputils

import (
	"fmt"
)

// api response general structure
type (
	APICode string

	Response struct {
		Message string   `json:"message"`
		Code    string   `json:"code"`
		Data    any      `json:"data"`
		Meta    any      `json:"meta,omitempty"`
		Errors  []string `json:"errors"`
	}
)

const (
	BAD_REQUEST           APICode = "bad_request"
	INTERNAL_SERVER       APICode = "internal_server"
	FORBIDDEN             APICode = "forbidden"
	SUCCESS               APICode = "success"
	UNAUTHORIZED          APICode = "unauthorized"
	TOKEN_REVOKED         APICode = "token_revoked"
	ACCESS_TOKEN_EXPIRED  APICode = "access_token_expired"
	REFRESH_TOKEN_EXPIRED APICode = "refresh_token_expired"
	// feel free to add more
)

var (
	Success                string = "success"
	ErrBadRequest          error  = fmt.Errorf("bad request")
	ErrForbidden           error  = fmt.Errorf("forbidden")
	ErrUnauthorized        error  = fmt.Errorf("unauthorized")
	ErrTokenRevoked        error  = fmt.Errorf("token revoked")
	ErrAccessTokenExpired  error  = fmt.Errorf("access token expired")
	ErrRefreshTokenExpired error  = fmt.Errorf("refresh token expired")
	// InternalServerErr error = fmt.Errorf("internal server error")
	// for internal server error, i think it's best to just use custom error instead of the pre-define one
)
