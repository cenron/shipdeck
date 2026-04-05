package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/cenron/shipdeck/internal/app"
	"github.com/cenron/shipdeck/internal/config"
	"github.com/cenron/shipdeck/internal/logging"
)

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log := logging.NewLogger()

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	log.Debug("ShipDeck is running...")

	err = app.Wire(ctx, cfg, log)
	if err != nil {
		return err
	}

	return nil
}
