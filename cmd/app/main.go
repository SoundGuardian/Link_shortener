package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"url-shortener/internal/config"
	"url-shortener/internal/http-server/handlers/redirect"
	"url-shortener/internal/http-server/handlers/remove"

	"url-shortener/internal/http-server/handlers/save"
	mw "url-shortener/internal/http-server/middleware"
	"url-shortener/internal/storage/postgre"
	"url-shortener/logger"
	"url-shortener/logger/sl"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)
	log = log.With(slog.String("env", cfg.Env))

	log.Info("initiolizing server", slog.String("address", cfg.Http.Address))
	log.Debug("logger debug mode enabled")

	storage, err := postgre.New(cfg.Data)
	if err != nil {
		log.Error("failed to initiolize storage", sl.Err(err))
	}
	fmt.Print(storage)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(mw.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			cfg.Http.User: cfg.Http.Password,
		}))

		r.Post("/", save.New(log, storage))
	})

	router.Delete("/{alias}", remove.New(log, storage))
	router.Get("/{alias}", redirect.New(log, storage))
	http.ListenAndServe(cfg.Http.Address, router)
	fmt.Println("HERE")
}
