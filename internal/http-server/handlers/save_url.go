package handlers

import (
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"

	"urlShortener/internal/http-server/response"
	"urlShortener/internal/storage"
	"urlShortener/internal/utils/random"
)

type Request struct {
	Url   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	response.Response
	Alias string `json:"alias"`
}

type UrlSaver interface {
	SaveUrl(urlToSave, alias string) (int64, error)
	AliasExists(alias string) (bool, error)
}

const aliasLength = 6

func SaveUrl(log *slog.Logger, urlSaver UrlSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.saveFile"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Error("empty request body")
				render.JSON(w, r, response.Error("empty request body"))

				return
			}
			log.Error("could not decode request", err)
			render.JSON(w, r, response.Error("could not decode request"))

			return
		}
		log.Info("request decoded", slog.Any("request", req))

		if err = validator.New().Struct(req); err != nil {
			var validationErr validator.ValidationErrors
			if errors.As(err, &validationErr) {
				log.Error("invalid request", slog.Any("errors", validationErr))
				render.JSON(w, r, response.ValidateError(validationErr))

				return
			} else {
				log.Error("could not validate request", err)
				render.JSON(w, r, response.Error("could not validate request"))

				return
			}
		}
		log.Info("validation passed")

		alias := req.Alias
		if alias == "" {
			var genErr error
			alias, genErr = GenerateAlias(urlSaver, aliasLength)

			if genErr != nil {
				log.Error("could not generate alias", genErr)
				render.JSON(w, r, response.Error("could not generate alias"))

				return
			}
			log.Info("alias generated", slog.String("alias", alias))
		}

		id, err := urlSaver.SaveUrl(req.Url, alias)
		if err != nil {
			if errors.Is(err, storage.ErrUrlExists) {
				log.Info("url already exists", slog.String("url", req.Url), slog.String("alias", alias))
				render.JSON(w, r, response.Error("alias already exists"))

				return
			}
			log.Error("could not save url", err)
			render.JSON(w, r, response.Error("could not save url"))

			return
		}
		log.Info("url saved", slog.Int64("id", id))

		render.JSON(w, r, Response{
			Response: response.Ok(),
			Alias:    alias,
		})
	}
}

func GenerateAlias(urlSaver UrlSaver, length int) (string, error) {
	for {
		alias := random.New(length)
		exists, err := urlSaver.AliasExists(alias)
		if err != nil {
			return "", err
		}
		if !exists {
			return alias, nil
		}
	}
}
