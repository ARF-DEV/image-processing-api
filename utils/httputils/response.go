package httputils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func SendResponse(w http.ResponseWriter, message string, data any, meta any, err error) {
	w.Header().Set("Content-Type", "application/json")

	// TODO handle error wrapper, code, etc
	errs := unwrapErrorStrs(err)
	code, apiCode := findCode(err)
	b := Response{
		Data:    data,
		Message: message,
		Meta:    meta,
		Code:    string(apiCode),
		Errors:  errs,
	}

	w.WriteHeader(code)
	jsonBody, _ := json.Marshal(b)
	fmt.Fprint(w, string(jsonBody))
}

func unwrapErrors(err error) []error {
	errs := []error{}
	for err != nil {
		errs = append(errs, err)
		err = errors.Unwrap(err)
	}

	return errs
}
func unwrapErrorStrs(err error) []string {
	errs := []string{}
	for err != nil {
		errs = append(errs, err.Error())
		err = errors.Unwrap(err)
	}

	return errs
}

func findCode(err error) (int, APICode) {
	errs := unwrapErrors(err)
	// Get first error
	if len(errs) == 0 {
		return http.StatusOK, SUCCESS
	}

	tarErr := errs[0]
	switch tarErr {
	case ErrBadRequest:
		return http.StatusBadRequest, BAD_REQUEST
	case ErrForbidden:
		return http.StatusForbidden, FORBIDDEN
	case ErrUnauthorized:
		return http.StatusUnauthorized, UNAUTHORIZED
	case ErrAccessTokenExpired:
		return http.StatusUnauthorized, ACCESS_TOKEN_EXPIRED
	case ErrRefreshTokenExpired:
		return http.StatusUnauthorized, REFRESH_TOKEN_EXPIRED
	case ErrTokenRevoked:
		return http.StatusUnauthorized, TOKEN_REVOKED
	default:
		return http.StatusInternalServerError, INTERNAL_SERVER
	}
}
