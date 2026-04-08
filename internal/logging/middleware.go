package logging

import (
	"fmt"
	"log/slog"
)

type WishMiddleware struct {
	Logger *slog.Logger
}

func (a WishMiddleware) Printf(format string, v ...interface{}) {
	a.Logger.Info(fmt.Sprintf(format, v...))
}
