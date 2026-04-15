package main

import "charm.land/lipgloss/v2"

type Theme struct {
	Title     lipgloss.Style
	ListTitle lipgloss.Style
	Selected  lipgloss.Style
	Inactive  lipgloss.Style
	Method    lipgloss.Style
}

func NewTheme() Theme {
	return Theme{
		// mauve title
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#cba6f7")).
			MarginBottom(1),

		ListTitle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#cba6f7")),

		// dark base text on blue background for current selection
		Selected: lipgloss.NewStyle().
			Background(lipgloss.Color("#89b4fa")).
			Foreground(lipgloss.Color("#1e1e2e")).
			PaddingLeft(1),

		// subtext1 for inactive items
		Inactive: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#a6adc8")),

		Method: lipgloss.NewStyle().
			Bold(true).
			Width(8),
	}
}
