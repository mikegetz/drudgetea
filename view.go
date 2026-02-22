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
	helpDescStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#4A4A4A"))
	helpKeyStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262"))
	helpSepStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#3C3C3C"))
	redStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff3c3c"))
	blueStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#5858fd"))
)

func (m model) View() string {
	centerStyle := lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center)

	view := "\n"

	view += m.TopHeadlineView()
	view += m.MainHeadlineView()
	view += centerStyle.Render(logo) + "\n"
	view += m.ColumnView()

	// The footer
	view += m.FooterView()

	// Debug info
	if os.Getenv("DEBUG") != "" {
		view += m.inputStyle.Render("Debug: "+fmt.Sprintf("cursorx: %d, cursory: %d, curMaxRow: %d, columnWidth: %d", m.cursorx, m.cursory, m.curMaxRow, m.columnWidth)) + "\n"
	}

	// The help view
	helpView := m.help.View(m.keys)
	height := 4 - strings.Count(helpView, "\n")

	view += strings.Repeat("\n", height) + helpView

	// Send the UI for rendering
	return view
}

func (m model) FooterView() string {
	view := ""
	view += "\n" + m.cursorStyle.Render(m.selected.Title)
	view += "\n(" + m.selected.Href + ")\n"
	view += helpDescStyle.Render("click to open")
	view += helpSepStyle.Render(" • ")
	view += helpKeyStyle.Render("c")
	view += helpDescStyle.Render(" copy link")
	return view
}

func (m model) TopHeadlineView() string {
	view := ""
	for _, mainHeadline := range m.topHeadlines {
		if mainHeadline.Color == "\033[31m" {
			view += redStyle.Align(lipgloss.Left).Render(mainHeadline.Title)
		} else {
			view += blueStyle.Align(lipgloss.Left).Render(mainHeadline.Title)
		}
		view += "\n"
	}
	return view
}

func (m model) MainHeadlineView() string {
	view := ""
	for _, mainHeadline := range m.mainHeadlines {
		if mainHeadline.Color == "\033[31m" {
			view += redStyle.Width(m.width).Align(lipgloss.Center).Render(mainHeadline.Title)
		} else {
			view += blueStyle.Width(m.width).Align(lipgloss.Center).Render(mainHeadline.Title)
		}
		view += "\n"
	}
	return view
}

func (m *model) ColumnView() string {
	rows := m.curMaxRow

	m.columnWidth = (m.width - (len(cs) * 2)) / len(m.headlines)
	view := ""
	if m.columnWidth > 3 {
		for i := 0; i < rows; i++ {
			for j, col := range m.headlines {
				if i < len(col) {
					var headlineStyle lipgloss.Style
					if col[i].Color == "\033[31m" {
						headlineStyle = redStyle.Width(m.columnWidth)
					} else {
						headlineStyle = blueStyle.Width(m.columnWidth)
					}
					if i == m.cursory && j == m.cursorx {
						m.cursorStyle = headlineStyle
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
