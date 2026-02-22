package main

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mikegetz/godrudge"
)

//go:embed logo
var logo string

const ansiReset = "\033[0m"
const ansiBold = "\033[1m"

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type keyMap struct {
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
	Help  key.Binding
	Quit  key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right}, // first column
		{k.Help, k.Quit},                // second column
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
}

type model struct {
	headlines  [][]godrudge.Headline // all headlines of drudge report
	cursorx    int                   // the current column
	cursory    int                   // the current row in the current column
	curMaxRow  int                   // the max number of rows in the current column
	width      int                   // width of the terminal
	maxRows    int                   // represents the column with the most headlines
	selected   godrudge.Headline     // the currently selected headline
	keys       keyMap
	help       help.Model
	inputStyle lipgloss.Style
	lastKey    string
	quitting   bool
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
		headlines: client.Page.HeadlineColumns,
		maxRows:   maxRows,
	}

	return model
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}
