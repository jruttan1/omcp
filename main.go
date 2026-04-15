package main

import (
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"fmt"
	"os"
)

func main() {
	ti := textinput.New()
	ti.Placeholder = "path/to/openapi.yaml"
	ti.Prompt = ""
	ti.SetWidth(46)
	ti.CharLimit = 200

	ui := textinput.New()
	ui.Placeholder = "https://..."
	ui.Prompt = ""
	ui.SetWidth(46)
	ui.CharLimit = 200

	s := spinner.New()
	s.Spinner = spinner.Dot

	d := delegate{selected: make(map[int]bool)}
	theme := NewTheme()
	l := list.New(nil, d, 0, 0)
	l.Title = "select endpoints to convert to mcp"
	l.Styles.Title = theme.ListTitle

	m := &model{
		state:    stateInput,
		input:    ti,
		urlInput: ui,
		spinner:  s,
		list:     l,
		delegate: d,
	}

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error launching CLI: %v", err)
		os.Exit(1)
	}
}
