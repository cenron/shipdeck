package session

import (
	"reflect"
	"testing"

	"github.com/charmbracelet/ssh"
)

type fakeSession struct {
	ssh.Session
	pty ssh.Pty
	ok  bool
}

func (f fakeSession) Pty() (ssh.Pty, <-chan ssh.Window, bool) {
	ch := make(chan ssh.Window)
	close(ch)
	return f.pty, ch, f.ok
}

func TestTeaHandlerUsesPTYWindowSize(t *testing.T) {
	t.Parallel()

	s := fakeSession{
		pty: ssh.Pty{Window: ssh.Window{Width: 120, Height: 40}},
		ok:  true,
	}

	model, opts := teaHandler(s)
	if len(opts) != 0 {
		t.Fatalf("expected no program options, got %d", len(opts))
	}

	v := reflect.ValueOf(model)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		t.Fatalf("expected struct model, got %s", v.Kind())
	}

	width := v.FieldByName("width")
	height := v.FieldByName("height")
	if !width.IsValid() || !height.IsValid() {
		t.Fatal("expected model to include width and height fields")
	}
	if width.Int() != 120 || height.Int() != 40 {
		t.Fatalf("expected width=120 height=40, got width=%d height=%d", width.Int(), height.Int())
	}
}
