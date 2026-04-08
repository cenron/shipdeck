package logging

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"
)

func TestWishMiddlewarePrintf(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := slog.New(&Handler{out: buf, level: slog.LevelInfo})

	m := WishMiddleware{Logger: logger}
	m.Printf("client %s connected", "alice")

	if got := buf.String(); !strings.Contains(got, "client alice connected") {
		t.Fatalf("expected formatted log line, got %q", got)
	}
	if !logger.Enabled(context.Background(), slog.LevelInfo) {
		t.Fatal("expected logger to be enabled for info")
	}
}
