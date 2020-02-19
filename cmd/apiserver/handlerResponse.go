package main

import (
	"net/http"

	"github.com/payfazz/go-errors"
	"github.com/payfazz/go-handler"
	"github.com/payfazz/go-handler/defresponse"

	"api/internal/storage"
)

func (h *httpHandler) resp404() *handler.Response {
	code := http.StatusNotFound

	data := struct {
		ErrorCode    int    `json:"error_code"`
		ErrorMessage string `json:"error_message"`
	}{
		ErrorCode:    code,
		ErrorMessage: http.StatusText(code),
	}

	return defresponse.JSON(code, data)
}

func (h *httpHandler) resp500(err error) *handler.Response {
	errors.PrintTo(h.errLogger, err)

	return defresponse.Status(http.StatusInternalServerError)
}

func (h *httpHandler) respMessage(r *http.Request, msg *storage.Message) *handler.Response {
	data := messageRepresentation{
		Title:   msg.Title,
		Message: msg.Message,
	}

	switch r.Header.Get("Accept") {
	case "application/x-www-form-urlencoded":
		return defresponse.Data(http.StatusOK, "application/x-www-form-urlencoded", data.asURLEncoded())
	case "application/json":
		fallthrough
	default:
		return defresponse.Data(http.StatusOK, "application/json", data.asJSON())
	}
}
