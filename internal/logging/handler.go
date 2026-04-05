package logging

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

type Handler struct {
	out   io.Writer
	level slog.Leveler
}

func NewLogger() *slog.Logger {
	return slog.New(&Handler{
		out:   os.Stdout,
		level: slog.LevelDebug,
	})
}

func (h *Handler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level.Level()
}

func (h *Handler) Handle(_ context.Context, r slog.Record) error {
	ts := r.Time.Format("2006-01-02 15:04:05")
	lvl := strings.ToUpper(r.Level.String())

	// Base format: DateTime [LogLevel => Message
	line := fmt.Sprintf("%s : [%s] => %s", ts, lvl, r.Message)

	// Optional: include structured attrs after the message.
	r.Attrs(func(a slog.Attr) bool {
		line += fmt.Sprintf(" %s=%v", a.Key, a.Value.Any())
		return true
	})

	line += "\n"
	_, err := io.WriteString(h.out, line)
	return err
}

func (h *Handler) WithAttrs(_ []slog.Attr) slog.Handler { return h }
func (h *Handler) WithGroup(_ string) slog.Handler      { return h }
