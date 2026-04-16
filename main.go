package main

import (
	"fmt"
	"os"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "init" {
		runInit()
		return
	}
	runServer()
}

func runInit() {
	ti := textinput.New()
	ti.Placeholder = "path/to/openapi.yaml"
	ti.Prompt = "> "
	ti.SetWidth(46)
	ti.CharLimit = 200

	ui := textinput.New()
	ui.Placeholder = "https://..."
	ui.Prompt = "> "
	ui.SetWidth(46)
	ui.CharLimit = 200

	ki := textinput.New()
	ki.Placeholder = "bearer token or api key"
	ki.Prompt = "> "
	ki.SetWidth(46)
	ki.CharLimit = 200

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
		keyInput: ki,
		spinner:  s,
		list:     l,
		delegate: d,
	}

	p := tea.NewProgram(m)
	result, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error launching CLI: %v\n", err)
		os.Exit(1)
	}
	m = result.(*model)

	cfg := Config{
		Spec:    m.input.Value(),
		BaseURL: m.urlInput.Value(),
		APIKey:  m.keyInput.Value(),
		Tools:   m.selectedTools,
		Info:    m.serverInfo,
	}
	if err := saveConfig(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "error saving config: %v\n", err)
		os.Exit(1)
	}

	path, _ := configPath()
	fmt.Fprintf(os.Stderr, "✓ config saved to %s\n\nAdd this to your MCP client config:\n\n{\n  \"mcpServers\": {\n    \"%s\": {\n      \"command\": \"omcp\"\n    }\n  }\n}\n", path, cfg.Info.Title)
}

func runServer() {
	cfg, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "no config found — run `omcp init` first\n")
		os.Exit(1)
	}

	mcpServer := createMcp(cfg.Spec, cfg.Info)
	createTools(cfg.Tools, mcpServer, cfg.BaseURL, cfg.APIKey)
	fmt.Fprintf(os.Stderr, "omcp: %s — %d tools registered on %s\n", cfg.Info.Title, len(cfg.Tools), cfg.BaseURL)

	if err := startServer(mcpServer); err != nil {
		fmt.Fprintf(os.Stderr, "error starting server: %v\n", err)
		os.Exit(1)
	}
}
