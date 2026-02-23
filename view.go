package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mikegetz/drudgetea/linkgloss"
)

var (
	cs             = "  "
	containerStyle = lipgloss.NewStyle().Padding(0, 2)
	borderColor    = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	redStyle       = linkgloss.New().Foreground(lipgloss.Color("#ff3c3c"))
	blueStyle      = linkgloss.New().Foreground(lipgloss.Color("#5858fd"))
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

	// The help view
	helpView := m.help.View(m.keys)
	height := 4 - strings.Count(helpView, "\n")

	view += strings.Repeat("\n", height) + helpView + "\n"

	// Send the UI for rendering
	return containerStyle.Render(view)
}

func (m model) FooterView() string {
	footerStyleColor := m.cursorStyle.UnsetAlign().UnsetWidth()
	footerStyleBorder := lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true).BorderForeground(lipgloss.Color("#888888")).Width((m.columnWidth * 3))

	// Debug info
	var debug string
	if os.Getenv("DEBUG") != "" {
		debug = m.inputStyle.Render("Debug: " + fmt.Sprintf("cursorx: %d, cursory: %d, curMaxRow: %d, columnWidth: %d", m.cursorx, m.cursory, m.curMaxRow, m.columnWidth))
	}

	view := "\n" + footerStyleBorder.Render(footerStyleColor.Render(m.selected.Title)+"\n("+m.selected.Href+")\n"+debug)

	return view
}

func (m *model) TopHeadlineView() string {
	view := ""
	for i, topHeadline := range m.topHeadlines {
		headlineStyle := m.TopHeadlineStyle(string(topHeadline.Color))

		if m.cursorGroup == 0 && i == m.cursory {
			m.cursorStyle = headlineStyle
			headlineStyle = headlineStyle.Bold(true)
		}

		view += headlineStyle.Href(topHeadline.Href).Render(topHeadline.Title)
		view += "\n"
	}
	return view
}

func (m *model) MainHeadlineView() string {
	view := ""
	for i, mainHeadline := range m.mainHeadlines {
		var headlineStyle = m.MainHeadlineStyle(string(mainHeadline.Color))

		if m.cursorGroup == 1 && i == m.cursory {
			m.cursorStyle = headlineStyle
			headlineStyle = headlineStyle.Bold(true)
		}

		view += headlineStyle.Href(mainHeadline.Href).Render(mainHeadline.Title)
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
					headlineStyle := m.HeadlineColumnStyle(string(col[i].Color))

					if i == m.cursory && j == m.cursorx && m.cursorGroup == 2 {
						m.cursorStyle = headlineStyle
						headlineStyle = headlineStyle.Bold(true)
					}

					if len(col[i].Title) > m.columnWidth {
						view += headlineStyle.Href(col[i].Href).Render(col[i].Title[:m.columnWidth-5] + "...")
					} else {
						view += headlineStyle.Href(col[i].Href).Render(col[i].Title)
					}
				} else {
					view += strings.Repeat(" ", m.columnWidth)
				}
				if j < len(m.headlines)-1 {
					view += cs
				}
			}
			view += "\n"
		}
	}
	return view
}

func (m model) HeadlineColumnStyle(color string) linkgloss.Style {
	var headlineStyle linkgloss.Style
	if color == "\033[31m" {
		headlineStyle = redStyle.Width(m.columnWidth)
	} else {
		headlineStyle = blueStyle.Width(m.columnWidth)
	}
	return headlineStyle
}

func (m model) MainHeadlineStyle(color string) linkgloss.Style {
	var headlineStyle linkgloss.Style
	if color == "\033[31m" {
		headlineStyle = redStyle.Width(m.width).Align(lipgloss.Center)
	} else {
		headlineStyle = blueStyle.Width(m.width).Align(lipgloss.Center)
	}
	return headlineStyle
}

func (m model) TopHeadlineStyle(color string) linkgloss.Style {
	var headlineStyle linkgloss.Style
	if color == "\033[31m" {
		headlineStyle = redStyle.Width(m.width).Align(lipgloss.Left)
	} else {
		headlineStyle = blueStyle.Width(m.width).Align(lipgloss.Left)
	}
	return headlineStyle
}
