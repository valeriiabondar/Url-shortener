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

type UrlGetter interface {
	GetUrl(alias string) (string, error)
}

func RedirectUrl(log *slog.Logger, urlGetter UrlGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.RedirectUrl"

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

		url, err := urlGetter.GetUrl(alias)
		if err != nil {
			if errors.Is(err, storage.ErrUrlNotFound) {
				log.Error("url not found", "alias", alias)
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, response.Error("url not found"))

				return
			}
			log.Error("could not get url", slog.Any("err", err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error("internal error"))

			return
		}

		log.Info("redirecting succeeded", slog.String("url", url), slog.String("alias", alias))

		http.Redirect(w, r, url, http.StatusFound)
	}
}
