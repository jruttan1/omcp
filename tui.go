package main

import (
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"fmt"
	"io"
	"strings"
)

type state int

const (
	stateInput state = iota
	stateLoading
	stateList
)

type toolsLoadedMsg []Tool
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
	"#cba6f7", "#c0b2f9", "#b4befe",
	"#9ec0fc", "#89b4fa", "#7dbef3",
	"#74c7ec", "#7fd8e8", "#89dceb",
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
	desc := fmt.Sprintf("   %s", tool.Summary)

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
	state        state
	input        textinput.Model
	urlInput     textinput.Model
	focusedField int // 0 = file path, 1 = base URL
	baseURL      string
	list         list.Model
	delegate     delegate
	spinner      spinner.Model
	err          error
	width        int
	height       int
}

func loadTools(fileName string) tea.Cmd {
	return func() tea.Msg {
		tools, err := parse(fileName)
		if err != nil {
			return toolsErrMsg(err)
		}
		return toolsLoadedMsg(tools)
	}
}

func (m *model) Init() tea.Cmd {
	return m.input.Focus()
}

func (m *model) focusField(i int) tea.Cmd {
	m.focusedField = i
	if i == 0 {
		m.urlInput.Blur()
		return m.input.Focus()
	}
	m.input.Blur()
	return m.urlInput.Focus()
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
		if m.state == stateInput {
			switch msg.String() {
			case "tab", "shift+tab":
				next := (m.focusedField + 1) % 2
				return m, m.focusField(next)
			case "enter":
				if m.focusedField == 0 {
					return m, m.focusField(1)
				}
				m.baseURL = m.urlInput.Value()
				m.state = stateLoading
				return m, tea.Batch(m.spinner.Tick, loadTools(m.input.Value()))
			}
		}
		if m.state == stateList && msg.String() == "space" && m.list.FilterState() != list.Filtering {
			i := m.list.Index()
			m.delegate.selected[i] = !m.delegate.selected[i]
			return m, nil
		}
	case toolsLoadedMsg:
		items := make([]list.Item, len(msg))
		for i, t := range msg {
			items[i] = t
		}
		m.list.SetItems(items)
		m.state = stateList
		return m, nil
	case toolsErrMsg:
		m.err = msg
		m.state = stateInput
		return m, nil
	}

	switch m.state {
	case stateInput:
		if m.focusedField == 0 {
			m.input, cmd = m.input.Update(msg)
		} else {
			m.urlInput, cmd = m.urlInput.Update(msg)
		}
	case stateLoading:
		m.spinner, cmd = m.spinner.Update(msg)
	case stateList:
		m.list, cmd = m.list.Update(msg)
	}

	return m, cmd
}

func (m *model) View() tea.View {
	var content string
	switch m.state {
	case stateInput:
		content = m.homepageView()
	case stateLoading:
		content = lipgloss.NewStyle().Padding(2, 4).Foreground(lipgloss.Color("#89b4fa")).Render(m.spinner.View() + " loading...")
	default:
		hint := lipgloss.NewStyle().Foreground(lipgloss.Color("#6c7086")).Render("space to toggle · enter to confirm")
		content = lipgloss.NewStyle().Padding(2, 4).Render(
			lipgloss.JoinVertical(lipgloss.Left, m.list.View(), hint),
		)
	}
	v := tea.NewView(content)
	v.AltScreen = true
	return v
}

const logo = `
_______/\\\\\_______/\\\\____________/\\\\________/\\\\\\\\\__/\\\\\\\\\\\\\___
 _____/\\\///\\\____\/\\\\\\________/\\\\\\_____/\\\////////__\/\\\/////////\\\_
  ___/\\\/__\///\\\__\/\\\//\\\____/\\\//\\\___/\\\/___________\/\\\_______\/\\\_
   __/\\\______\//\\\_\/\\\\///\\\/\\\/_\/\\\__/\\\_____________\/\\\\\\\\\\\\\/__
    _\/\\\_______\/\\\_\/\\\__\///\\\/___\/\\\_\/\\\_____________\/\\\/////////____
     _\//\\\______/\\\__\/\\\____\///_____\/\\\_\//\\\____________\/\\\_____________
      __\///\\\__/\\\____\/\\\_____________\/\\\__\///\\\__________\/\\\_____________
       ____\///\\\\\/_____\/\\\_____________\/\\\____\////\\\\\\\\\_\/\\\_____________
        ______\/////_______\///______________\///________\/////////__\///______________`

func (m *model) homepageView() string {
	dim := lipgloss.NewStyle().Foreground(lipgloss.Color("#bac2de"))

	logoLines := strings.Split(logo, "\n")
	for i, line := range logoLines {
		logoLines[i] = lipgloss.NewStyle().Foreground(lipgloss.Color(logoColors[i%len(logoColors)])).Render(line)
	}
	title := strings.Join(logoLines, "\n")

	errStr := ""
	if m.err != nil {
		errStr = "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("#f38ba8")).Render("✗ "+m.err.Error())
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

	return lipgloss.NewStyle().Padding(2, 4).Render(lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		"",
		"",
		dim.Render("turn an openapi spec into an mcp server in 1 command :)"),
		"",
		dim.Render("api spec file"),
		box(m.input, m.focusedField == 0)+errStr,
		"",
		dim.Render("server base url"),
		box(m.urlInput, m.focusedField == 1),
		"",
		dim.Render("tab to switch · enter to confirm · ctrl+c to quit"),
	))
}
