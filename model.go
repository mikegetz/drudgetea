package main

import (
	_ "embed"
	"fmt"
	"os"
	"strings"
	"time"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/lipgloss/v2"
	tea "charm.land/bubbletea/v2"
	"github.com/mikegetz/godrudge"
)

//go:embed logo
var logo string

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type keyMap struct {
	Up            key.Binding
	Down          key.Binding
	Left          key.Binding
	Right         key.Binding
	Help          key.Binding
	Quit          key.Binding
	Copy          key.Binding
	ShowLogo      key.Binding
	Less          key.Binding
	Tab           key.Binding
	Version       key.Binding
	DisableLinks  key.Binding
	ToggleRefresh key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Tab, k.Less, k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},                   // first column
		{k.Less, k.Tab, k.ShowLogo, k.Copy},               // second column
		{k.Help, k.Quit, k.DisableLinks, k.ToggleRefresh}, // third column
		{k.Version},
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move right"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Copy: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "copy url"),
	),
	Less: key.NewBinding(
		key.WithKeys("enter", "space"),
		key.WithHelp("space/enter", "toggle more rows"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "toggle headline group"),
	),
	ShowLogo: key.NewBinding(
		key.WithKeys("L"),
		key.WithHelp("L", "toggle ascii art logo"),
	),
	Version: key.NewBinding(
		key.WithKeys("v"),
		key.WithHelp("", strings.Repeat("\n", 3)+"Version: "+Version),
	),
	DisableLinks: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "toggle links"),
	),
	ToggleRefresh: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "toggle refresh"),
	),
}

type model struct {
	// data
	client *godrudge.Client
	time   time.Time // last time of RSS feed fetch

	// view state
	cursorGroup      int            // the current column group (top, main, or headline columns)
	cursorx          int            // the current column
	cursory          int            // the current row in the current column
	curMaxRow        int            // the max number of rows in the current column
	width            int            // width of the terminal
	maxRows          int            // represents the column with the most headlines
	columnWidth      int            // the width of each column
	toggleRowLess    int            // toggle to expand column rows, value represents the max rows when toggled
	showLogo         bool           // whether to show the ascii art logo
	cursorStyle      lipgloss.Style //current cursor style - remove when godrudge supports lipgloss styles
	disableLinkgloss bool           // whether to disable linkgloss styles, for better compatibility with terminals that don't support them

	//controller state
	selected       godrudge.Headline // the currently selected headline
	keys           keyMap            // the keybindings
	help           help.Model        // the help view
	inputStyle     lipgloss.Style    // style for debug info
	lastKey        string            // the last key pressed, for debug purposes
	quitting       bool              // whether the application is quitting
	refreshEnabled bool              // whether auto-refresh is enabled
}

func initialModel() model {
	client := godrudge.NewClient()
	err := client.ParseRSS()
	if err != nil {
		fmt.Printf("Error parsing drudgereport: %v\n", err)
		os.Exit(1)
	}

	maxRows := refreshMaxRows(client)

	h := help.New()
	h.Styles = help.DefaultDarkStyles()

	model := model{
		client:         client,
		time:           time.Now(),
		toggleRowLess:  10,
		cursorGroup:    1,
		showLogo:       true,
		maxRows:        maxRows,
		keys:           keys,
		help:           h,
		inputStyle:     lipgloss.NewStyle().Foreground(lipgloss.Color("#FF75B7")),
		refreshEnabled: true,
	}

	return model
}

func refreshMaxRows(client *godrudge.Client) int {
	maxRows := 0
	for _, col := range client.Page.HeadlineColumns {
		if len(col) > maxRows {
			maxRows = len(col)
		}
	}
	return maxRows
}

func (m model) Init() tea.Cmd {
	return tea.Batch(refresh(30*time.Second), tea.RequestBackgroundColor)
}

type tickMsg time.Time

func refresh(d time.Duration) tea.Cmd {
	// tea.Tick schedules a single message after d.
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
