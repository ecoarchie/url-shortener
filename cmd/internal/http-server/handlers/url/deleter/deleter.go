package deleter

import (
	"log/slog"
	"net/http"

	resp "github.com/ecoarchie/url-shortener/cmd/internal/lib/api/response"
	"github.com/ecoarchie/url-shortener/cmd/internal/lib/logger/slg"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type URLDeleter interface {
	DeleteURL(alias string) error
}

func New(logger *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"

		logger = logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			logger.Info("alias is empty")
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}

		err := urlDeleter.DeleteURL(alias)

		if err != nil {
			logger.Error("unable to delete url", slg.Err(err))
			render.JSON(w, r, resp.Error("internal error"))
			return
		}

		logger.Info("url deleted")
		render.JSON(w, r, resp.OK())
	}
}