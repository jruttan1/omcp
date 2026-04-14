package main

import "charm.land/lipgloss/v2"

type Theme struct {
	Title    lipgloss.Style
	Selected lipgloss.Style
	Inactive lipgloss.Style
	Method   lipgloss.Style
}

func NewTheme() Theme {
	return Theme{
		// purple title
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("99")).
			MarginBottom(1),

		// white text on a blue background for the current selection
		Selected: lipgloss.NewStyle().
			Background(lipgloss.Color("62")).
			Foreground(lipgloss.Color("255")).
			PaddingLeft(1),

		// dimmer gray for the other items
		Inactive: lipgloss.NewStyle().
			Foreground(lipgloss.Color("245")),

		Method: lipgloss.NewStyle().
			Bold(true).
			Width(8),
	}
}
