package ui

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestNewAndInit(t *testing.T) {
	m := New(10, 20)
	if m == nil {
		t.Fatal("expected model")
	}
	if m.width != 10 || m.height != 20 {
		t.Fatalf("unexpected size: %dx%d", m.width, m.height)
	}
	if m.ready {
		t.Fatal("expected ready false")
	}
	if m.status != "connected" {
		t.Fatalf("unexpected status: %q", m.status)
	}
	if cmd := m.Init(); cmd != nil {
		t.Fatal("expected nil init command")
	}
}

func TestUpdateWindowSize(t *testing.T) {
	m := New(0, 0)
	model, cmd := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	if cmd != nil {
		t.Fatal("expected nil command")
	}
	if model != m {
		t.Fatal("expected same model pointer")
	}
	if m.width != 80 || m.height != 24 || !m.ready {
		t.Fatalf("unexpected model state: %+v", m)
	}
}

func TestUpdateQuitKeys(t *testing.T) {
	m := New(0, 0)

	_, cmd := m.Update(tea.KeyPressMsg(tea.Key{Text: "q", Code: 'q'}))
	if cmd == nil {
		t.Fatal("expected quit command for q")
	}

	_, cmd = m.Update(tea.KeyPressMsg(tea.Key{Code: 'c', Mod: tea.ModCtrl}))
	if cmd == nil {
		t.Fatal("expected quit command for ctrl+c")
	}
}

func TestView(t *testing.T) {
	m := New(0, 0)
	loading := m.View().Content
	if !strings.Contains(strings.ToLower(loading), "loading") {
		t.Fatalf("expected loading view, got %q", loading)
	}

	m.ready = true
	v := m.View().Content
	if !strings.Contains(v, "ShipDeck") || !strings.Contains(v, "Projects") || !strings.Contains(v, "q quit") {
		t.Fatalf("unexpected dashboard view: %q", v)
	}
}
