package main

import (
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"fmt"
)

type model struct {
	tools  []Tool
	cursor int
	styles Theme
}

// on app startup
func (m *model) Init() tea.Cmd {
	return nil
}

// handle inputs and update state
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	keyMsg, ok := msg.(tea.KeyMsg) // convert msg to tea.KeyMsg, return result and success
	if !ok {
		return m, nil
	}

	switch keyMsg.String() { // unpack KeyMsg interface to a string
	case "up":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down":
		if m.cursor < len(m.tools)-1 {
			m.cursor++
		}
	case "q", "ctrl+c":
		return m, tea.Quit
	}

	return m, nil
}

func (m *model) View() tea.View {

	lines := make([]string, len(m.tools))

	for i, tool := range m.tools {
		content := fmt.Sprintf("%s %s", tool.Method, tool.Route)
		if m.cursor == i {
			lines[i] = m.styles.Selected.Render("> " + content)
		} else {
			lines[i] = m.styles.Inactive.Render(" " + content)

		}
	}

	return tea.NewView(lipgloss.JoinVertical(lipgloss.Left, lines...))
}
