package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"
)

const gray = "\033[90m"
const reset = "\033[0m"

type MultilineHandler struct {
	io.Writer
}

func (h *MultilineHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return true
}

func (h *MultilineHandler) Handle(_ context.Context, r slog.Record) error {
	fmt.Fprintf(h.Writer, "%s [%s] %s \n", r.Time.Format(time.RFC3339), r.Level, r.Message)

	r.Attrs(func(a slog.Attr) bool {
		h.printAttr(a)
		return true
	})

	return nil
}

func (h *MultilineHandler) printAttr(a slog.Attr) {
	if a.Value.Kind() == slog.KindGroup {
		for _, sub := range a.Value.Group() {
			fmt.Fprintf(h.Writer, "%s%s.%s: %v%s\n", gray, a.Key, sub.Key, sub.Value, reset)
		}
	} else {
		fmt.Fprintf(h.Writer, "  %s%s%s: %v\n", gray, a.Key, reset, a.Value)
	}
}

func (h *MultilineHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *MultilineHandler) WithGroup(_ string) slog.Handler {
	return h
}
