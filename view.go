package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	cs := " | "
	columnSeparator := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).Render(cs)
	centerStyle := lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center)
	redStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
	blueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#3030fe"))

	view := centerStyle.Render(logo)

	view += "\n"

	maxColumnSize := 0
	for _, col := range m.headlines {
		if len(col) > maxColumnSize {
			maxColumnSize = len(col)
		}
	}
	columnWidth := (m.width - (len(cs) * 2)) / len(m.headlines)

	if columnWidth > 3 {
		for i := 0; i < maxColumnSize; i++ {
			for j, col := range m.headlines {
				if i < len(col) {
					var headlineStyle lipgloss.Style
					if col[i].Color == "\033[31m" {
						headlineStyle = redStyle.Width(columnWidth)
					} else {
						headlineStyle = blueStyle.Width(columnWidth)
					}
					if i == m.cursory && j == m.cursorx {
						m.selected = col[i]
						headlineStyle = headlineStyle.Bold(true)
					}

					if len(col[i].Title) > columnWidth {
						view += headlineStyle.Render(col[i].Title[:columnWidth-3] + "...")
					} else {
						view += headlineStyle.Render(col[i].Title)
					}
				} else {
					view += strings.Repeat(" ", columnWidth)
				}
				if j < len(m.headlines)-1 {
					view += columnSeparator
				}
			}
			view += "\n"
		}
	}

	// The footer

	view += "\n Selected: " + m.selected.Title + "\n(" + m.selected.Href + ")\n"

	if os.Getenv("DEBUG") != "" {
		view += m.inputStyle.Render("Debug: "+fmt.Sprintf("cursorx: %d, cursory: %d, curMaxRow: %d, columnWidth: %d", m.cursorx, m.cursory, m.curMaxRow, columnWidth)) + "\n"
	}

	// The help view
	helpView := m.help.View(m.keys)
	height := 3 - strings.Count(helpView, "\n")

	view += strings.Repeat("\n", height) + helpView

	// Send the UI for rendering
	return view
}
