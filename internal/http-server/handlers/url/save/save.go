package save

import (
	"errors"
	"log/slog"
	"net/http"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/slogx"
	"url-shortener/internal/lib/random"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

// TODO: maybe move to config
const generatedAliasLength = 8

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			msg := "failed to decode request body"
			log.Error(msg, slogx.Err(err))

			render.JSON(w, r, resp.Error(msg))

			return
		}

		log.Info("request body decode successfully", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			response := resp.ValidationError(validateErr)

			log.Error(response.Error, slogx.Err(err))

			render.JSON(w, r, response)
			return
		}

		alias := req.Alias

		if alias == "" {
			alias, err = random.NewRandomString(generatedAliasLength)

			if err != nil {
				msg := "failed to generate unique alias"
				log.Error(msg, slogx.Err(err))

				render.JSON(w, r, resp.Error(msg))

				return
			}
		}

		id, err := urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			msg := "url already exists"
			log.Info(msg, slog.String("url", req.URL))
			render.JSON(w, r, resp.Error(msg))
			return
		}
		if err != nil {
			msg := "failed to add url"
			log.Error(msg, slogx.Err(err))
			render.JSON(w, r, resp.Error(msg))
			return
		}

		log.Info("url added", slog.Int64("id", id))

		render.JSON(w, r, Response{Response: resp.OK(), Alias: alias})
	}
}
