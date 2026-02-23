package main

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mikegetz/drudgetea/linkgloss"
	"github.com/mikegetz/godrudge"
)

//go:embed logo
var logo string

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type keyMap struct {
	Up       key.Binding
	Down     key.Binding
	Left     key.Binding
	Right    key.Binding
	Help     key.Binding
	Quit     key.Binding
	Copy     key.Binding
	ShowLogo key.Binding
	Less     key.Binding
	Tab      key.Binding
	Version  key.Binding
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
		{k.Up, k.Down, k.Left, k.Right},     // first column
		{k.Less, k.Tab, k.ShowLogo, k.Copy}, // second column
		{k.Help, k.Quit, keys.Version},      // third column
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
		key.WithKeys("enter", " "),
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
		key.WithHelp("", "\nVersion: "+Version),
	),
}

type model struct {
	// data
	headlines     [][]godrudge.Headline // all headlines of drudge report
	mainHeadlines []godrudge.Headline   // the main headlines, which are displayed above the logo section
	topHeadlines  []godrudge.Headline   // the top headlines, which are displayed above the main headlines left aligned

	// view state
	cursorGroup   int             // the current column group (top, main, or headline columns)
	cursorx       int             // the current column
	cursory       int             // the current row in the current column
	curMaxRow     int             // the max number of rows in the current column
	width         int             // width of the terminal
	maxRows       int             // represents the column with the most headlines
	columnWidth   int             // the width of each column
	toggleRowLess int             // toggle to expand column rows, value represents the max rows when toggled
	showLogo      bool            // whether to show the ascii art logo
	cursorStyle   linkgloss.Style //current cursor style - remove when godrudge supports lipgloss styles

	//controller state
	selected   godrudge.Headline // the currently selected headline
	keys       keyMap            // the keybindings
	help       help.Model        // the help view
	inputStyle lipgloss.Style    // style for debug info
	lastKey    string            // the last key pressed, for debug purposes
	quitting   bool              // whether the application is quitting
}

func initialModel() model {
	client := godrudge.NewClient()
	err := client.ParseRSS()
	if err != nil {
		fmt.Printf("Error parsing drudgereport: %v\n", err)
		os.Exit(1)
	}

	maxRows := 0
	for _, col := range client.Page.HeadlineColumns {
		if len(col) > maxRows {
			maxRows = len(col)
		}
	}

	model := model{
		headlines:     client.Page.HeadlineColumns,
		mainHeadlines: client.Page.MainHeadlines,
		topHeadlines:  client.Page.TopHeadlines,
		toggleRowLess: 10,
		cursorGroup:   1,
		showLogo:      true,
		maxRows:       maxRows,
		keys:          keys,
		help:          help.New(),
		inputStyle:    lipgloss.NewStyle().Foreground(lipgloss.Color("#FF75B7")),
	}

	return model
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}
