package main

import (
	"fmt"
	"io"
	"net/url"
	"strings"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
)

// This file is kinda hard to read because of the nested switches for state, view, input and different structs and functions interacting
// Im gonna come back and refactor this into functions for each view state eventually

type state int

const (
	stateInput state = iota
	stateLoading
	stateList
)

type toolsLoadedMsg struct {
	tools []Tool
	info  Info
}
type toolsErrMsg error

// delegate renders each list row and tracks which items are selected.
// because maps are reference types, m.delegate and the copy inside m.list
// share the same underlying data, toggling m.delegate.selected is selects item
type delegate struct {
	selected map[int]bool
}

func (d delegate) Height() int                             { return 2 }
func (d delegate) Spacing() int                            { return 0 }
func (d delegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

const checkPrefixLen = 3

var logoColors = []string{
	"#f5c2e7", "#cba6f7", "#b4befe",
	"#89b4fa", "#74c7ec", "#89dceb",
	"#74c7ec", "#89b4fa", "#b4befe",
}

func (d delegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	tool, ok := item.(Tool)
	if !ok {
		return
	}

	check := "○"
	if d.selected[index] {
		check = "●"
	}

	// ○  method /route
	// matched indices are into FilterValue, so offset by checkPrefixLen
	title := fmt.Sprintf("%s  %s %s", check, tool.Method, tool.Route)
	desc := fmt.Sprintf("%s", tool.Summary)

	isSelected := index == m.Index()
	isFiltered := m.FilterState() == list.Filtering || m.FilterState() == list.FilterApplied

	if isFiltered {
		matches := m.MatchesForItem(index)
		offset := make([]int, len(matches))
		for i, r := range matches {
			offset[i] = r + checkPrefixLen
		}
		unmatched := lipgloss.NewStyle().Inline(true)
		matched := unmatched.Underline(true)
		title = lipgloss.StyleRunes(title, offset, matched, unmatched)
	}

	if isSelected {
		fmt.Fprintf(w, "%s\n%s",
			lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#cba6f7")).Render(title),
			lipgloss.NewStyle().Foreground(lipgloss.Color("#a6adc8")).Render(desc),
		)
	} else {
		fmt.Fprintf(w, "%s\n%s",
			title,
			lipgloss.NewStyle().Foreground(lipgloss.Color("#6c7086")).Render(desc),
		)
	}
}

type model struct {
	state         state
	input         textinput.Model
	urlInput      textinput.Model
	keyInput      textinput.Model
	focusedField  int // 0 = file path, 1 = base URL, 2 = key
	baseURL       string
	serverInfo    Info
	list          list.Model
	delegate      delegate
	spinner       spinner.Model
	err           error
	urlErr        error
	selectedTools []Tool
	width         int
	height        int
}

func loadTools(fileName string) tea.Cmd {
	return func() tea.Msg {
		tools, info, err := parse(fileName)
		if err != nil {
			return toolsErrMsg(err)
		}
		return toolsLoadedMsg{tools: tools, info: info}
	}
}

func (m *model) Init() tea.Cmd {
	return m.input.Focus()
}

func (m *model) focusField(i int) tea.Cmd {
	m.focusedField = i
	switch i {
	case 0:
		m.keyInput.Blur()
		m.urlInput.Blur()
		return m.input.Focus()
	case 1:
		m.keyInput.Blur()
		m.input.Blur()
		return m.urlInput.Focus()
	case 2:
		m.input.Blur()
		m.urlInput.Blur()
		return m.keyInput.Focus()
	}
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.list.SetSize(msg.Width-8, msg.Height-6)
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		switch m.state { // switch based on state of ui
		case stateInput: // initial input screen
			switch msg.String() { // switch based on user command for input screen
			case "tab", "shift+tab": // tab to switch fields
				next := (m.focusedField + 1) % 3
				return m, m.focusField(next)
			case "enter": // enter switches or continues to next screen if both boxes are full
				if m.focusedField == 1 {
					u := m.urlInput.Value()
					parsed, parseErr := url.Parse(u)
					if parseErr != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
						m.urlErr = fmt.Errorf("invalid url: must start with http:// or https://")
						return m, nil
					}
					m.urlErr = nil
				}
				if m.focusedField < 2 {
					return m, m.focusField(m.focusedField + 1)
				}
				m.baseURL = m.urlInput.Value()
				m.state = stateLoading
				return m, tea.Batch(m.spinner.Tick, loadTools(m.input.Value()))
			default:
				switch m.focusedField {
				case 0:
					m.input, cmd = m.input.Update(msg)
				case 1:
					m.urlInput, cmd = m.urlInput.Update(msg)
				case 2:
					m.keyInput, cmd = m.keyInput.Update(msg)
				}
			}
		case stateList: // second screen for viewing enpoints and server data
			switch msg.String() {
			case "space": // space to select an endpoint
				if m.list.FilterState() != list.Filtering {
					i := m.list.Index()
					m.delegate.selected[i] = !m.delegate.selected[i]
					return m, nil
				}
				m.list, cmd = m.list.Update(msg)
			case "enter": // confirm selection
				if m.list.FilterState() == list.Filtering {
					m.list, cmd = m.list.Update(msg)
					return m, cmd
				}
				selected := false
				for i, v := range m.delegate.selected { // m.delegate.selected is a map[int]bool so check if there is at least 1 bool == true
					if v {
						selected = true
						m.selectedTools = append(m.selectedTools, m.list.Items()[i].(Tool)) // complex because it unpacks list, gets item at index i and asserts to Tool type
					}
				}
				if !selected {
					m.err = fmt.Errorf("select at least one endpoint to continue")
				} else {
					return m, tea.Quit // successfully entered selections, quit tui, server start handled in main
				}

			default:
				m.list, cmd = m.list.Update(msg)
			}
		case stateLoading: // between screens render spinner
			m.spinner, cmd = m.spinner.Update(msg)
		}

	case toolsLoadedMsg: // async success state from parser to ui
		m.serverInfo = msg.info
		items := make([]list.Item, len(msg.tools))
		for i, t := range msg.tools {
			items[i] = t
		}
		m.list.SetItems(items)
		m.state = stateList
		return m, nil

	case toolsErrMsg: // async fail state from parser to ui
		m.err = msg
		m.state = stateInput
		return m, nil

	default:
		if m.state == stateList {
			m.list, cmd = m.list.Update(msg)
		}
	}

	return m, cmd
}

func (m *model) View() tea.View {
	var content string
	switch m.state { // view states corresponding to update switch
	case stateInput:
		content = m.homepageView()
	case stateLoading:
		content = lipgloss.NewStyle().Padding(2, 4).Foreground(lipgloss.Color("#89b4fa")).Render(m.spinner.View() + " loading...")
	default:
		title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#cba6f7")).Render(m.serverInfo.Title)
		meta := lipgloss.NewStyle().Foreground(lipgloss.Color("#6c7086")).Render("v" + m.serverInfo.Version)
		desc := lipgloss.NewStyle().Foreground(lipgloss.Color("#a6adc8")).Render(m.serverInfo.Description)
		hint := lipgloss.NewStyle().Foreground(lipgloss.Color("#6c7086")).Render("space to toggle · enter to confirm")
		errStr := ""
		if m.err != nil {
			errStr = "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("#f38ba8")).Render("✗ "+m.err.Error())
		}
		content = lipgloss.NewStyle().Padding(2, 4).Render(
			lipgloss.JoinVertical(lipgloss.Left, title+" "+meta, desc, "", m.list.View(), hint, errStr),
		)
	}
	v := tea.NewView(content)
	v.AltScreen = true
	return v
}

const logo = `
 ██████╗ ███╗   ███╗ ██████╗██████╗ 
██╔═══██╗████╗ ████║██╔════╝██╔══██╗
██║   ██║██╔████╔██║██║     ██████╔╝
██║   ██║██║╚██╔╝██║██║     ██╔═══╝ 
╚██████╔╝██║ ╚═╝ ██║╚██████╗██║     
 ╚═════╝ ╚═╝     ╚═╝ ╚═════╝╚═╝ `

func (m *model) homepageView() string {
	dim := lipgloss.NewStyle().Foreground(lipgloss.Color("#6c7086"))
	normal := lipgloss.NewStyle().Foreground(lipgloss.Color("#bac2de"))
	accent := lipgloss.NewStyle().Foreground(lipgloss.Color("#89b4fa"))

	label := func(text string, fieldIndex int) string {
		if m.focusedField == fieldIndex {
			return accent.Render(text)
		}
		return dim.Render(text)
	}

	logoLines := strings.Split(logo, "\n")
	for i, line := range logoLines {
		logoLines[i] = lipgloss.NewStyle().Foreground(lipgloss.Color(logoColors[i%len(logoColors)])).Render(line)
	}
	title := strings.Join(logoLines, "\n")

	errStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#f38ba8"))
	errStr := ""
	if m.err != nil {
		errStr = "\n" + errStyle.Render("✗ "+m.err.Error())
	}
	urlErrStr := ""
	if m.urlErr != nil {
		urlErrStr = "\n" + errStyle.Render("✗ "+m.urlErr.Error())
	}

	box := func(input textinput.Model, focused bool) string {
		borderColor := lipgloss.Color("#313244")
		if focused {
			borderColor = lipgloss.Color("#89b4fa")
		}
		return lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Padding(0, 1).Width(50).
			Render(input.View())
	}

	arrow := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#cba6f7")). // catppuccin mauve
		Render("→")

	return lipgloss.NewStyle().Padding(2, 4).Render(lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		normal.Render("openapi spec ")+arrow+normal.Render(" mcp server in 1 command :)"),
		"",
		label("api spec file", 0),
		box(m.input, m.focusedField == 0)+errStr,
		"",
		label("base url of a server", 1),
		box(m.urlInput, m.focusedField == 1)+urlErrStr,
		"",
		label("api key (if auth is needed)", 2),
		box(m.keyInput, m.focusedField == 2),
		dim.Render("tab to switch · enter to confirm · ctrl+c to quit"),
	))
}
