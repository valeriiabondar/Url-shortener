package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"urlShortener/internal/http-server/response"
	"urlShortener/internal/storage"
)

type UrlDeleter interface {
	DeleteUrl(alias string) error
}

func DeleteUrl(log *slog.Logger, urlDeleter UrlDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.DeleteUrl"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")

		if strings.TrimSpace(alias) == "" {
			log.Error("empty alias")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid request"))
			return
		}

		err := urlDeleter.DeleteUrl(alias)
		if err != nil {
			if errors.Is(err, storage.ErrUrlNotFound) {
				log.Error("url not found", "alias", alias)
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, response.Error("url not found"))

				return
			}
			log.Error("could not delete url", slog.Any("err", err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error("internal error"))

			return
		}

		log.Info("url deleted", slog.String("alias", alias))

		w.WriteHeader(http.StatusNoContent)
	}
}
