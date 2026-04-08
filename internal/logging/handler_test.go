package logging

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"strings"
	"testing"
	"time"
)

func TestHandlerEnabled(t *testing.T) {
	h := &Handler{level: slog.LevelInfo}
	if h.Enabled(context.Background(), slog.LevelDebug) {
		t.Fatal("expected debug to be disabled")
	}
	if !h.Enabled(context.Background(), slog.LevelInfo) {
		t.Fatal("expected info to be enabled")
	}
}

func TestNewLogger(t *testing.T) {
	if got := NewLogger(); got == nil {
		t.Fatal("expected logger instance")
	}
}

func TestHandlerHandleWritesFormattedLine(t *testing.T) {
	buf := &bytes.Buffer{}
	h := &Handler{out: buf, level: slog.LevelDebug}

	r := slog.NewRecord(time.Date(2026, 4, 8, 9, 30, 0, 0, time.UTC), slog.LevelWarn, "hello", 0)
	r.AddAttrs(slog.String("k", "v"))

	if err := h.Handle(context.Background(), r); err != nil {
		t.Fatalf("Handle() returned error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "[WARN]") {
		t.Fatalf("expected WARN level in output, got %q", got)
	}
	if !strings.Contains(got, "=> hello") {
		t.Fatalf("expected message in output, got %q", got)
	}
	if !strings.Contains(got, "k=v") {
		t.Fatalf("expected attrs in output, got %q", got)
	}
}

func TestHandlerHandleWriteError(t *testing.T) {
	h := &Handler{out: errWriter{}, level: slog.LevelDebug}
	r := slog.NewRecord(time.Now(), slog.LevelInfo, "x", 0)

	err := h.Handle(context.Background(), r)
	if err == nil {
		t.Fatal("expected write error")
	}
}

func TestWithAttrsAndWithGroupReturnHandler(t *testing.T) {
	h := &Handler{}
	if got := h.WithAttrs([]slog.Attr{slog.String("x", "y")}); got != h {
		t.Fatal("WithAttrs should return same handler")
	}
	if got := h.WithGroup("g"); got != h {
		t.Fatal("WithGroup should return same handler")
	}
}

type errWriter struct{}

func (errWriter) Write(_ []byte) (int, error) {
	return 0, errors.New("boom")
}
