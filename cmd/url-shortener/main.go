package main

import (
    "log/slog"
    "net/http"
    "os"
    "os/signal"
    "syscall"

    "github.com/DimsFromDergachy/Url-Shortener/internal/config"
    "github.com/DimsFromDergachy/Url-Shortener/internal/http-server/handlers/url/redirect"
    "github.com/DimsFromDergachy/Url-Shortener/internal/http-server/handlers/url/save"
    mwLogger "github.com/DimsFromDergachy/Url-Shortener/internal/http-server/middleware/logger"
    "github.com/DimsFromDergachy/Url-Shortener/internal/lib/logger/sl"
    "github.com/DimsFromDergachy/Url-Shortener/internal/storage/sqlite"
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
)

const (
    envLocal = "local"
    envDev = "dev"
    envProd = "prod"
)

func main() {
    cfg := config.MustLoad()

    log := setupLogger(cfg.Env)
    log = log.With(slog.String("env", cfg.Env))

    log.Info("initializing server", slog.String("address", cfg.Address))
    log.Debug("logger debug mode enabled")

    storage, err := sqlite.New(cfg.StoragePath)
    if err != nil {
        log.Error("failed to initialize storage", sl.Err(err))
    }

    router := chi.NewRouter()

    router.Use(middleware.RequestID)
    router.Use(middleware.Logger)
    router.Use(middleware.Recoverer)
    router.Use(middleware.URLFormat)
    router.Use(mwLogger.New(log))
    router.Post("/", save.New(log, storage))
    router.Get("/{alias}", redirect.New(log, storage))

    log.Info("starting server")

    done := make(chan os.Signal, 1)
    signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

    srv := &http.Server {
        Addr:         cfg.Address,
        Handler:      router,
        ReadTimeout:  cfg.HTTPServer.Timeout,
        WriteTimeout: cfg.HTTPServer.Timeout,
        IdleTimeout:  cfg.HTTPServer.IdleTimeout,
    }

    go func() {
        if err := srv.ListenAndServe(); err != nil {
            log.Error("failed to start server", sl.Err(err))
        }
    }()

    log.Info("server started")

    <- done

    log.Info("stopping server")
}

func setupLogger(env string) *slog.Logger {
    var log *slog.Logger

    switch env {
    case envLocal:
        log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
    case envDev:
        log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
    case envProd:
        log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
    }

    return log
}