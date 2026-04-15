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

	// run tui
	p := tea.NewProgram(m)
	result, err := p.Run() // p.Run() returns the state of the tui run
	if err != nil {
		fmt.Printf("Error launching CLI: %v", err)
		os.Exit(1)
	}
	m = result.(*model) // cast the tea.Model that Run() returns to custom model pointer

	// start mcp
	fmt.Print("Starting MCP Server...")
	mcpServer := createMcp(m.input.Value(), m.serverInfo)
	createTools(m.selectedTools, mcpServer, m.urlInput.Value())

	if err := startServer(mcpServer); err != nil {
		fmt.Printf("Error occured while starting server: %v\n", err)
		os.Exit(1)
	}

}
