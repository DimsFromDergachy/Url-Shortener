package save

import (
    "errors"
    "io"
    "log/slog"
    "net/http"

    resp "github.com/DimsFromDergachy/Url-Shortener/internal/lib/api/response"
    "github.com/DimsFromDergachy/Url-Shortener/internal/lib/logger/sl"
    "github.com/DimsFromDergachy/Url-Shortener/internal/lib/random"
    "github.com/DimsFromDergachy/Url-Shortener/internal/storage"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/go-chi/render"
    "github.com/go-playground/validator"
)

type Request struct {
    URL   string `json:"url" validate:"required,url"`
    Alias string `json:"alias,omitempty"`
}

type Response struct {
    resp.Response
    Alias  string `json:"alias"`
}

const aliasLength = 6

type URLSaver interface {
    SaveURL(URL, alias string) (int64, error)
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        const op = "handlers.url.save.New"

        log = log.With(
            slog.String("op", op),
            slog.String("request_id", middleware.GetReqID(r.Context())),
        )

        var req Request

        err := render.DecodeJSON(r.Body, &req)
        if errors.Is(err, io.EOF) {
            log.Error("request body is empty")

            render.JSON(w, r, resp.Error("empty request"))

            return
        }
        if err != nil {
            log.Error("failed to decode request body", sl.Err(err))

            render.JSON(w, r, resp.Error("failed to decode request"))

            return
        }

        log.Info("request body decoded", slog.Any("req", req))

        if err := validator.New().Struct(req); err != nil {
            validateErr := err.(validator.ValidationErrors)

            log.Error("invalid request", sl.Err(err))

            render.JSON(w, r, resp.ValidationError(validateErr))

            return
        }

        alias := req.Alias
        if alias == "" {
            alias = random.NewRandomString(aliasLength)
        }

        id, err := urlSaver.SaveURL(req.URL, alias)
        if errors.Is(err, storage.ErrURLExists) {
            log.Info("url already exists", slog.String("url", req.URL))

            render.JSON(w, r, resp.Error("URL already exists"))

            return
        }
        if err != nil {
            log.Error("failed to add url", sl.Err(err))

            render.JSON(w, r, resp.Error("failed to add URL"))

            return
        }

        log.Info("url added", slog.Int64("id", id))

        responseOK(w, r, alias)
    }
}

func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
    render.JSON(w, r, Response{
        Response: resp.OK(),
        Alias: alias,
    })
}