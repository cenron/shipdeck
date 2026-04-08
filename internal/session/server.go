package session

import (
	"context"
	"errors"
	"fmt"
	"net"
	"path/filepath"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/log/v2"
	"charm.land/wish/v2"
	"charm.land/wish/v2/activeterm"
	"charm.land/wish/v2/bubbletea"
	"charm.land/wish/v2/logging"
	"github.com/cenron/shipdeck/internal/config"
	"github.com/cenron/shipdeck/internal/ui"
	"github.com/charmbracelet/ssh"
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
		wish.WithAddress(net.JoinHostPort(s.config.Host, s.config.Port)),
		wish.WithAuthorizedKeys(authKeys),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			activeterm.Middleware(), // Bubble Tea apps usually require a PTY.
			logging.Middleware(),
		),
	)

	if err != nil {
		return fmt.Errorf("create SSH server: %w", err)
	}

	errChan := make(chan error, 1)
	log.Info("Starting SSH server", "host", s.config.Host, "port", s.config.Port)
	go func() {
		if listenErr := serv.ListenAndServe(); listenErr != nil && !errors.Is(listenErr, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", listenErr)
			errChan <- listenErr
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

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, _ := s.Pty()
	m := ui.New(pty.Window.Width, pty.Window.Height)
	return m, []tea.ProgramOption{}
}
