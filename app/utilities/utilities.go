package utilities

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
)

type ErrorResponse struct {
	StatusText string `json:"status"`
	AppCode    int64  `json:"code,omitempty"`
	ErrorText  string `json:"error"`
}

func RenderHTTPError(w http.ResponseWriter, r *http.Request) {
	jsonError(w, r, 400, errors.New(http.StatusText(400)))
}

func RenderAuthError(w http.ResponseWriter, r *http.Request) {
	jsonError(w, r, 401, errors.New(http.StatusText(401)))
}

func RenderNotAllowedError(w http.ResponseWriter, r *http.Request) {
	jsonError(w, r, 405, errors.New(http.StatusText(405)))
}

func RenderUnprocessableError(w http.ResponseWriter, r *http.Request) {
	jsonError(w, r, 422, errors.New(http.StatusText(422)))
}

func RenderServerError(w http.ResponseWriter, r *http.Request, err error) {
	jsonError(w, r, 500, err)
}

func jsonError(w http.ResponseWriter, r *http.Request, code int, err error) {
	w.WriteHeader(code)

	render.JSON(w, r, ErrorResponse{
		StatusText: http.StatusText(code),
		ErrorText:  err.Error(),
	})
}
