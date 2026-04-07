package session

import (
	"context"
	"errors"
	"fmt"
	"net"
	"path/filepath"
	"time"

	"charm.land/log/v2"
	"charm.land/wish/v2"
	"github.com/cenron/shipdeck/internal/config"
	"github.com/charmbracelet/ssh"
)

const (
	host = "localhost"
	port = "23234"
)

type Server struct {
	ctx    context.Context
	config config.Config
}

func NewServer(ctx context.Context, cfg config.Config) *Server {
	return &Server{
		ctx:    ctx,
		config: cfg,
	}
}

func (s *Server) Run() error {
	authKeys, err := filepath.Abs(s.config.AuthorizedKeysPath)
	if err != nil {
		return fmt.Errorf("resolve authorized keys path: %w", err)
	}

	serv, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithAuthorizedKeys(authKeys),
	)

	if err != nil {
		return fmt.Errorf("create SSH server: %w", err)
	}

	errChan := make(chan error, 1)
	log.Info("Starting SSH server", "host", host, "port", port)
	go func() {
		if err = serv.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", err)
			errChan <- err
		}
	}()

	select {
	case <-s.ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		return serv.Shutdown(shutdownCtx)
	case err := <-errChan:
		return err
	}
}
