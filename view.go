package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mikegetz/drudgetea/linkgloss"
)

var (
	cs              = "  "
	columnSeparator = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).Render(cs)
	helpDescStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#4A4A4A"))
	helpKeyStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262"))
	helpSepStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#3C3C3C"))
	redStyle        = linkgloss.New().Foreground(lipgloss.Color("#ff3c3c"))
	blueStyle       = linkgloss.New().Foreground(lipgloss.Color("#5858fd"))
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
	view += "\n" + m.cursorStyle.UnsetAlign().UnsetWidth().Render(m.selected.Title)
	view += "\n(" + m.selected.Href + ")\n"
	view += helpDescStyle.Render("click to open")
	view += helpSepStyle.Render(" • ")
	view += helpKeyStyle.Render("c")
	view += helpDescStyle.Render(" copy link")
	return view
}

func (m *model) TopHeadlineView() string {
	view := ""
	for i, mainHeadline := range m.topHeadlines {
		var headlineStyle linkgloss.Style
		if mainHeadline.Color == "\033[31m" {
			headlineStyle = redStyle.Align(lipgloss.Left).Width(m.width)
		} else {
			headlineStyle = blueStyle.Align(lipgloss.Left).Width(m.width)
		}

		if m.cursorGroup == 0 && i == m.cursory {
			m.cursorStyle = headlineStyle
			headlineStyle = headlineStyle.Bold(true)
		}

		view += headlineStyle.Render(mainHeadline.Title)
		view += "\n"
	}
	return view
}

func (m *model) MainHeadlineView() string {
	view := ""
	for i, mainHeadline := range m.mainHeadlines {
		var headlineStyle linkgloss.Style
		if mainHeadline.Color == "\033[31m" {
			headlineStyle = redStyle.Width(m.width).Align(lipgloss.Center)
		} else {
			headlineStyle = blueStyle.Width(m.width).Align(lipgloss.Center)
		}

		if m.cursorGroup == 1 && i == m.cursory {
			m.cursorStyle = headlineStyle
			headlineStyle = headlineStyle.Bold(true)
		}

		view += headlineStyle.Render(mainHeadline.Title)
		view += "\n"
	}
	return view
}

func (m *model) ColumnView() string {
	var rows int
	if m.toggleRowLess > 0 {
		rows = m.toggleRowLess
	} else {
		rows = m.maxRows
	}

	m.columnWidth = (m.width - (len(cs) * 2)) / len(m.headlines)
	view := ""
	if m.columnWidth > 3 {
		for i := 0; i < rows; i++ {
			for j, col := range m.headlines {
				if i < len(col) {
					var headlineStyle linkgloss.Style
					if col[i].Color == "\033[31m" {
						headlineStyle = redStyle.Width(m.columnWidth)
					} else {
						headlineStyle = blueStyle.Width(m.columnWidth)
					}

					if i == m.cursory && j == m.cursorx && m.cursorGroup == 2 {
						m.cursorStyle = headlineStyle
						headlineStyle = headlineStyle.Bold(true)
					}

					if len(col[i].Title) > m.columnWidth {
						view += headlineStyle.Href(col[i].Href).Render(col[i].Title[:m.columnWidth-3] + "...")
					} else {
						view += headlineStyle.Href(col[i].Href).Render(col[i].Title)
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
