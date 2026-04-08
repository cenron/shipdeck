package session

import (
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/wish/v2"
	"charm.land/wish/v2/activeterm"
	"charm.land/wish/v2/bubbletea"
	"charm.land/wish/v2/testsession"
	"github.com/charmbracelet/ssh"
	gossh "golang.org/x/crypto/ssh"
)

func TestAuthorizedKeyStartsBubbleTeaSession(t *testing.T) {
	authorizedSigner := mustGenerateSigner(t)
	authorizedKeysPath := writeAuthorizedKeysFile(t, authorizedSigner)

	server, err := wish.NewServer(
		wish.WithAddress("127.0.0.1:0"),
		wish.WithAuthorizedKeys(authorizedKeysPath),
		wish.WithMiddleware(
			bubbletea.Middleware(func(_ ssh.Session) (tea.Model, []tea.ProgramOption) {
				return quitModel{}, nil
			}),
			activeterm.Middleware(),
		),
	)
	if err != nil {
		t.Fatalf("create wish server: %v", err)
	}

	addr := testsession.Listen(t, server)

	if err := connectAndRunShell(addr, authorizedSigner); err != nil {
		t.Fatalf("expected authorized key to start Bubble Tea session: %v", err)
	}
}

func TestUnauthorizedKeyIsRejected(t *testing.T) {
	authorizedSigner := mustGenerateSigner(t)
	unauthorizedSigner := mustGenerateSigner(t)
	authorizedKeysPath := writeAuthorizedKeysFile(t, authorizedSigner)

	server, err := wish.NewServer(
		wish.WithAddress("127.0.0.1:0"),
		wish.WithAuthorizedKeys(authorizedKeysPath),
		wish.WithMiddleware(
			bubbletea.Middleware(func(_ ssh.Session) (tea.Model, []tea.ProgramOption) {
				return quitModel{}, nil
			}),
			activeterm.Middleware(),
		),
	)
	if err != nil {
		t.Fatalf("create wish server: %v", err)
	}

	addr := testsession.Listen(t, server)

	err = connectAndRunShell(addr, unauthorizedSigner)
	if err == nil {
		t.Fatal("expected unauthorized key to be rejected")
	}
}

type quitModel struct{}

func (quitModel) Init() tea.Cmd { return tea.Quit }

func (m quitModel) Update(tea.Msg) (tea.Model, tea.Cmd) { return m, nil }

func (quitModel) View() tea.View { return tea.NewView("ok") }

func mustGenerateSigner(t *testing.T) gossh.Signer {
	t.Helper()
	_, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	signer, err := gossh.NewSignerFromKey(priv)
	if err != nil {
		t.Fatalf("new signer: %v", err)
	}
	return signer
}

func writeAuthorizedKeysFile(t *testing.T, signer gossh.Signer) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "authorized_keys")
	line := gossh.MarshalAuthorizedKey(signer.PublicKey())
	if err := os.WriteFile(path, line, 0o600); err != nil {
		t.Fatalf("write authorized_keys: %v", err)
	}
	return path
}

func connectAndRunShell(addr string, signer gossh.Signer) error {
	client, err := gossh.Dial("tcp", addr, &gossh.ClientConfig{
		User:            "test",
		Auth:            []gossh.AuthMethod{gossh.PublicKeys(signer)},
		HostKeyCallback: gossh.InsecureIgnoreHostKey(),
		Timeout:         2 * time.Second,
	})
	if err != nil {
		return err
	}
	defer func() { _ = client.Close() }()

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer func() { _ = session.Close() }()

	if err := session.RequestPty("xterm", 24, 80, gossh.TerminalModes{}); err != nil {
		return err
	}
	if err := session.Shell(); err != nil {
		return err
	}

	waitCh := make(chan error, 1)
	go func() {
		waitCh <- session.Wait()
	}()

	select {
	case err := <-waitCh:
		return err
	case <-time.After(2 * time.Second):
		return fmt.Errorf("session did not exit")
	}
}
