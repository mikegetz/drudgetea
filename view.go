package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	cs              = " | "
	columnSeparator = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).Render(cs)
	redStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
	blueStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#3030fe"))
)

func (m model) View() string {
	centerStyle := lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center)
	view := centerStyle.Render(logo)

	view += "\n"

	view += m.ColumnView()

	// The footer

	view += "\n Selected: " + m.selected.Title + "\n(" + m.selected.Href + ")\n"

	if os.Getenv("DEBUG") != "" {
		view += m.inputStyle.Render("Debug: "+fmt.Sprintf("cursorx: %d, cursory: %d, curMaxRow: %d, columnWidth: %d", m.cursorx, m.cursory, m.curMaxRow, m.columnWidth)) + "\n"
	}

	// The help view
	helpView := m.help.View(m.keys)
	height := 3 - strings.Count(helpView, "\n")

	view += strings.Repeat("\n", height) + helpView

	// Send the UI for rendering
	return view
}

func (m *model) ColumnView() string {
	m.columnWidth = (m.width - (len(cs) * 2)) / len(m.headlines)
	view := ""
	if m.columnWidth > 3 {
		for i := 0; i < m.maxRows; i++ {
			for j, col := range m.headlines {
				if i < len(col) {
					var headlineStyle lipgloss.Style
					if col[i].Color == "\033[31m" {
						headlineStyle = redStyle.Width(m.columnWidth)
					} else {
						headlineStyle = blueStyle.Width(m.columnWidth)
					}
					if i == m.cursory && j == m.cursorx {
						m.selected = col[i]
						headlineStyle = headlineStyle.Bold(true)
					}

					if len(col[i].Title) > m.columnWidth {
						view += headlineStyle.Render(col[i].Title[:m.columnWidth-3] + "...")
					} else {
						view += headlineStyle.Render(col[i].Title)
					}
				} else {
					view += strings.Repeat(" ", m.columnWidth)
				}
				if j < len(m.headlines)-1 {
					view += columnSeparator
				}
			}
			view += "\n"
		}
	}
	return view
}
