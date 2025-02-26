package save

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"url-shortener/internal/config"
	"url-shortener/internal/lib/random"
	"url-shortener/internal/storage"
	"url-shortener/logger/sl"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
)

type URLSaver interface {
	SaveURL(URL, alias string) (int64, error)
}
type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
	Alias  string `json:"alias"`
}

const (
	StatusOk    = "OK"
	StatusError = "Error"
)

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.NEW"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			render.JSON(w, r,
				Error("empty request"))
			return
		}
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, Error("failed to decode body"))

			return
		}

		log.Info("request cody decoded", slog.Any("req", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, Error(validateErr.Error()))
			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(config.AliasLen)
		}
		id, err := urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExist) {
			log.Info("url already exists", slog.String("url", req.URL))

			render.JSON(w, r, Error("url already exists"))
			return
		}
		if err != nil {
			log.Error("failed to add url", sl.Err(err))

			render.JSON(w, r, Error("failed to add url"))
			return
		}

		log.Info("url added", slog.Int64("id", id))

		responseOK(w, r, alias)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Status: StatusOk,
		Alias:  alias,
	})
}
