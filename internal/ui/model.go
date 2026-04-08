package ui

import tea "charm.land/bubbletea/v2"

type Model struct {
	width  int
	height int
	ready  bool
	status string
	err    error
}

func New(width, height int) *Model {
	return &Model{
		width:  width,
		height: height,
		ready:  false,
		status: "connected",
		err:    nil,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}
