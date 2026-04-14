package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	tools  []Tool
	cursor int
}

// on app startup
func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// input, return and updating state

	switch msg := msg.(type)


	return m, nil
}

func (m model) View() string {
	// display
	return ""
}
