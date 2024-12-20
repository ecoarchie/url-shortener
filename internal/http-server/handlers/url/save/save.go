package save

import (
	"errors"
	"log/slog"
	"net/http"

	resp "github.com/ecoarchie/url-shortener/internal/lib/api/response"
	"github.com/ecoarchie/url-shortener/internal/lib/logger/slg"
	"github.com/ecoarchie/url-shortener/internal/lib/random"
	"github.com/ecoarchie/url-shortener/internal/storage"
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
	SaveURL(urlToSave, alias string) (int64, error)
}

// TODO maybe move to config or to db
const aliasLength = 4

func New(logger *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		logger = logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			logger.Error("failed to decode request body", slg.Err(err))
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		logger.Info("req body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			logger.Error("invalid request", slg.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))
			return
		}
		alias := req.Alias
		if alias == "" {
			alias = random.RandomString(aliasLength)
		}

		id, err := urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			logger.Info("url already exists", slog.String("url", req.URL))
			render.JSON(w, r, resp.Error("url already exists"))
			return
		}

		logger.Info("url added", slog.Int64("id", id))

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Alias:    alias,
		})
	}
}
