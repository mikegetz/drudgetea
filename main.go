package main

import (
	_ "embed"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mikegetz/godrudge"
)

//go:embed logo
var logo string

const ansiReset = "\033[0m"
const ansiBold = "\033[1m"

type model struct {
	headlines [][]godrudge.Headline // all headlines of drudge report
	cursorx   int                   // the current column
	cursory   int                   // the current row in the current column
	curMaxRow int                   // the max number of rows in the current column
	width     int                   // width of the terminal
	maxRows   int                   // represents the column with the most headlines
	selected  godrudge.Headline     // the currently selected headline
}

func initialModel() model {
	client := godrudge.NewClient()
	err := client.ParseRSS()
	if err != nil {
		fmt.Printf("Error parsing drudgereport: %v\n", err)
		os.Exit(1)
	}
	model := model{
		headlines: client.Page.HeadlineColumns,
		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
	}

	maxRows := 0
	for _, col := range model.headlines {
		if len(col) > maxRows {
			maxRows = len(col)
		}
	}

	model.maxRows = maxRows

	return model
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		//The "left" and "h" keys move the cursor left
		case "left", "h":
			if m.cursorx > 0 {
				if (len(m.headlines[m.cursorx-1]) - 1) >= m.cursory {
					m.cursorx--
				}
			}

		// The "right" and "l" keys move the cursor right
		case "right", "l":
			if m.cursorx < len(m.headlines)-1 {
				if (len(m.headlines[m.cursorx+1]) - 1) >= m.cursory {
					m.cursorx++
				}
			}

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursory > 0 {
				m.cursory--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursory < m.curMaxRow-1 {
				m.cursory++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			//TODO: load page?
		}

	}

	m.curMaxRow = len(m.headlines[m.cursorx])
	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	line := logo + "\n"

	maxColumnSize := 0
	for _, col := range m.headlines {
		if len(col) > maxColumnSize {
			maxColumnSize = len(col)
		}
	}
	columnWidth := m.width / len(m.headlines)

	if columnWidth > 3 {
		for i := 0; i < maxColumnSize; i++ {
			for j, col := range m.headlines {
				if i < len(col) {
					if i == m.cursory && j == m.cursorx {
						m.selected = col[i]
						line += ansiBold
					}

					if len(col[i].Title) > columnWidth {
						line += fmt.Sprintf(string(col[i].Color)+"%-*.*s"+ansiReset, columnWidth, columnWidth, col[i].Title[:columnWidth-3]+"...")
					} else {
						line += fmt.Sprintf(string(col[i].Color)+"%-*s"+ansiReset, columnWidth, col[i].Title)
					}
				} else {
					line += fmt.Sprintf("%-*s", columnWidth, "")
				}
				if j < len(m.headlines)-1 {
					line += " | "
				}
			}
			line += "\n"
		}
	}

	// The footer

	line += "\n Selected: " + m.selected.Title + "\n(" + m.selected.Href + ")\n"

	if len(os.Args) > 1 && os.Args[1] == "--debug" {
		line += "Debug: " + fmt.Sprintf("cursorx: %d, cursory: %d, curMaxRow: %d", m.cursorx, m.cursory, m.curMaxRow) + "\n"
	}

	line += "\nPress q to quit.\n"

	// Send the UI for rendering
	return line
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("error starting bubbletea: %v\n", err)
		os.Exit(1)
	}
}
