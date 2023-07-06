package handler

import (
	"net/http"

	"github.com/go-chi/render"
)

// ErrResponse renderer type for handling all sorts of errors.
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

// Render sets the application-specific error code in AppCode.
func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// ErrInvalidRequest returns status and error message.
func ErrInvalidRequest(err error, httpStatusCode int) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: httpStatusCode,
		StatusText:     http.StatusText(httpStatusCode),
		ErrorText:      err.Error(),
	}
}
