package remove

import (
	"errors"
	"log/slog"
	"net/http"
	"url-shortener/internal/http-server/handlers/save"
	"url-shortener/internal/storage"
	"url-shortener/logger/sl"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type URLDeleter interface {
	DeleteURL(alias string) (string, error)
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.delete.delete.go"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")

		if alias == "" {
			log.Info("alias is empty")

			render.JSON(w, r, save.Error("not found"))
			return
		}

		delStat, err := urlDeleter.DeleteURL(alias)

		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", "alias", alias)

			render.JSON(w, r, save.Error("not found"))
			return
		}
		if err != nil {
			log.Error("failed to delete url", sl.Err(err))

			render.JSON(w, r, save.Error("internal error"))

			return
		}
		log.Info("url sucessfully deleted", slog.String("url", delStat))
	}
}
