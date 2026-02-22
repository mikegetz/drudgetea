package main

import (
	"fmt"
	"os"
)

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
