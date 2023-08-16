package slogdiscard

import (
    "context"
    "log/slog"
)

func NewDiscardLogger() *slog.Logger {
    return slog.New(NewDiscardHandler())
}

type DiscardHandler struct{}

func NewDiscardHandler() *DiscardHandler {
    return &DiscardHandler{}
}

func (h *DiscardHandler) Handle(_ context.Context, _ slog.Record) error {
    return nil
}

func (h *DiscardHandler) Enabled(context.Context, slog.Level) bool {
    return false
}

func (h *DiscardHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
    return h
}

func (h *DiscardHandler) WithGroup(name string) slog.Handler {
    return h
}

