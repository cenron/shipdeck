package ui

import tea "charm.land/bubbletea/v2"

func (m *Model) View() tea.View {
	if !m.ready {
		return tea.NewView("“loading terminal ...")
	}

	header := "ShipDeck"
	body := "Projects: --\nDeployments: --\nAlerts: --"
	footer := "q quit"
	screen := header + "\n\n" + body + "\n\n" + footer
	
	return tea.NewView(screen)
}
