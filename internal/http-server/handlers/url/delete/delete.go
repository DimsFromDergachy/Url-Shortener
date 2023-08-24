package delete

import (
    "log/slog"
    "net/http"

    resp "github.com/DimsFromDergachy/Url-Shortener/internal/lib/api/response"
    "github.com/DimsFromDergachy/Url-Shortener/internal/lib/logger/sl"
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/go-chi/render"
)

//go:generate go run github.com/vektra/mockery/v2@v2.32.4 --name=URLRemover
type URLRemover interface {
    DeleteURL(alias string) error
}

func New(log *slog.Logger, urlRemover URLRemover) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        const op = "handlers.url.delete.New"

        log = log.With(
            slog.String("op", op),
            slog.String("request_id", middleware.GetReqID(r.Context())),
        )

        alias := chi.URLParam(r, "alias")
        if alias == "" {
            log.Info("alias is empty")

            render.JSON(w, r, resp.Error("alias is empty"))

            return
        }

        err := urlRemover.DeleteURL(alias)
        if err != nil {
            log.Error("failed to remove url", sl.Err(err))

            render.JSON(w, r, resp.Error("internal error"))

            return
        }

        log.Info("removed url")

        render.JSON(w, r, resp.OK())
    }
}
