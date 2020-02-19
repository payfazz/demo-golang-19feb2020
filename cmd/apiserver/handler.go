package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/payfazz/go-errors"
	"github.com/payfazz/go-handler"
	"github.com/payfazz/go-handler/defresponse"
	"github.com/payfazz/go-middleware"
	"github.com/payfazz/go-router/defhandler"
	"github.com/payfazz/go-router/method"
	"github.com/payfazz/go-router/path"
	"github.com/payfazz/go-router/segment"
	"github.com/payfazz/stdlog"

	"api/internal/storage"
)

type httpHandler struct {
	storage   storage.Storage
	errLogger stdlog.Printer
}

func createHandler(storage storage.Storage, errLogger stdlog.Printer) http.HandlerFunc {
	h := &httpHandler{
		storage:   storage,
		errLogger: errLogger,
	}

	return middleware.Compile(
		h.mustValidAuth(),

		path.H{
			"/:id": middleware.Compile(
				segment.MustEnd,
				method.H{
					"GET":    h.get(),
					"PUT":    h.createOrUpdate(),
					"DELETE": h.delete(),
				}.C(),
			),
		}.C(),
	)
}

func (h *httpHandler) mustValidAuth() func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			user, pass, _ := r.BasicAuth()

			if user == "admin" && pass == "password" {
				next(w, r)
				return
			}

			defhandler.StatusUnauthorized(w, r)
		}
	}
}

func (h *httpHandler) get() http.HandlerFunc {
	return handler.Of(func(r *http.Request) *handler.Response {
		id, ok := segment.Get(r, "id")
		if !ok {
			return h.resp404()
		}

		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return h.resp404()
		}

		msg, err := h.storage.GetMessage(r.Context(), storage.MessageID(int(idInt)))
		if err != nil {
			return h.resp500(errors.Wrap(err))
		}

		if msg == nil {
			return h.resp404()
		}

		return h.respMessage(r, msg)
	})
}

func (h *httpHandler) createOrUpdate() http.HandlerFunc {
	return handler.Of(func(r *http.Request) *handler.Response {
		id, ok := segment.Get(r, "id")
		if !ok {
			return defresponse.Status(http.StatusBadRequest)
		}

		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return defresponse.Status(http.StatusBadRequest)
		}

		data, err := ioutil.ReadAll(io.LimitReader(r.Body, 1<<20+1)) // 1MiB
		if err != nil {
			return defresponse.Status(http.StatusBadRequest)
		}
		if len(data) >= 1<<20 {
			return defresponse.Status(http.StatusRequestEntityTooLarge)
		}

		var msg messageRepresentation

		switch r.Header.Get("Content-Type") {
		case "application/json":
			err = msg.parseFromJSON(data)
			if err != nil {
				return defresponse.Status(http.StatusBadRequest)
			}
		case "application/x-www-form-urlencoded":
			err = msg.parseFromURLEncoded(data)
			if err != nil {
				return defresponse.Status(http.StatusBadRequest)
			}
		default:
			return defresponse.Status(http.StatusUnsupportedMediaType)
		}

		storageMsg := &storage.Message{
			ID:      storage.MessageID(idInt),
			Title:   msg.Title,
			Message: msg.Message,
		}
		err = h.storage.StoreMessage(r.Context(), storageMsg)
		if err != nil {
			return h.resp500(errors.Wrap(err))
		}

		return h.respMessage(r, storageMsg)
	})
}

func (h *httpHandler) delete() http.HandlerFunc {
	return handler.Of(func(r *http.Request) *handler.Response {
		id, ok := segment.Get(r, "id")
		if !ok {
			return h.resp404()
		}

		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return h.resp404()
		}

		err = h.storage.DeleteMessage(r.Context(), storage.MessageID(int(idInt)))
		if err != nil {
			return h.resp500(errors.Wrap(err))
		}

		return defresponse.Status(http.StatusOK)
	})
}
