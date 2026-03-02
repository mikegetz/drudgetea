package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/mikegetz/godrudge"
)

var (
	cs             = "  "
	containerStyle = lipgloss.NewStyle().Padding(0, 2)
)

func (m model) View() tea.View {
	centerStyle := lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center)

	view := "\n"

	view += m.TopHeadlineView()
	view += m.MainHeadlineView()

	if m.showLogo {
		view += centerStyle.Render(logo) + "\n"
	} else {
		view += "\n"
	}

	view += m.ColumnView()

	// The footer
	view += m.FooterView()

	// The help view
	m.help.SetWidth(m.width - 4)
	helpView := m.help.View(m.keys)
	height := 4 - strings.Count(helpView, "\n")

	view += strings.Repeat("\n", height) + helpView + "\n"

	// Send the UI for rendering

	teaView := tea.NewView(containerStyle.Render(view))

	teaView.BackgroundColor = lipgloss.Color("#181A1B")

	return teaView
}

func (m model) FooterView() string {
	footerStyleColor := m.cursorStyle.UnsetAlign().UnsetWidth()
	footerStyleBorder := lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true).BorderForeground(lipgloss.Color("#888888")).Width((m.columnWidth * 3))

	// Debug info
	var debug string
	if os.Getenv("DEBUG") != "" {
		debug = m.inputStyle.Render("\nDebug: " + fmt.Sprintf("cursorx: %d, cursory: %d, curMaxRow: %d, columnWidth: %d, disableLinkgloss: %t, time: %s", m.cursorx, m.cursory, m.curMaxRow, m.columnWidth, m.disableLinkgloss, m.time.Format("15:04:05")))
	}

	view := footerStyleColor.Render(m.selected.Title)

	urlLine := "(" + lipgloss.NewStyle().Hyperlink(m.selected.URL).Render(m.selected.URL) + ")"
	timeSince := time.Since(m.time).Truncate(time.Second).String()
	timeLabel := m.help.Styles.ShortKey.Render(timeSince) + m.help.Styles.ShortDesc.Render(" since last refresh")
	remainingWidth := m.columnWidth*3 - 2 - lipgloss.Width(urlLine)
	timeStr := lipgloss.NewStyle().Width(remainingWidth).Align(lipgloss.Right).Render(timeLabel)
	view += "\n" + urlLine + timeStr

	view += debug

	return footerStyleBorder.Render(view)
}

func (m *model) TopHeadlineView() string {
	var topHeadlines []godrudge.Headline
	if m.toggleRowLess > 0 && m.toggleRowLess < len(m.client.Page.TopHeadlines) {
		topHeadlines = m.client.Page.TopHeadlines[:m.toggleRowLess]
	} else {
		topHeadlines = m.client.Page.TopHeadlines
	}
	view := ""
	for i, topHeadline := range topHeadlines {
		headlineStyle := topHeadline.Style.Width(m.width).Align(lipgloss.Left).Hyperlink(topHeadline.URL)

		if m.disableLinkgloss {
			headlineStyle = headlineStyle.UnsetHyperlink()
		}

		if m.cursorGroup == 0 && i == m.cursory {
			m.cursorStyle = headlineStyle
			headlineStyle = headlineStyle.Bold(true)
		}

		view += headlineStyle.Render(topHeadline.Title)
		view += "\n"
	}
	view += "\n"
	return view
}

func (m *model) MainHeadlineView() string {
	var mainHeadlines []godrudge.Headline
	if m.toggleRowLess > 0 && m.toggleRowLess < len(m.client.Page.MainHeadlines) {
		mainHeadlines = m.client.Page.MainHeadlines[:m.toggleRowLess]
	} else {
		mainHeadlines = m.client.Page.MainHeadlines
	}
	view := ""
	for i, mainHeadline := range mainHeadlines {
		var headlineStyle = mainHeadline.Style.Width(m.width).Align(lipgloss.Center).Hyperlink(mainHeadline.URL)

		if m.disableLinkgloss {
			headlineStyle = headlineStyle.UnsetHyperlink()
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

	m.columnWidth = (m.width - (len(cs) * 2)) / len(m.client.Page.HeadlineColumns)
	view := ""
	if m.columnWidth > 3 {
		for i := 0; i < rows; i++ {
			for j, col := range m.client.Page.HeadlineColumns {
				if i < len(col) {
					headlineStyle := col[i].Style.Width(m.columnWidth).Hyperlink(col[i].URL)

					if m.disableLinkgloss {
						headlineStyle = headlineStyle.UnsetHyperlink()
					}

					if i == m.cursory && j == m.cursorx && m.cursorGroup == 2 {
						m.cursorStyle = headlineStyle
						headlineStyle = headlineStyle.Bold(true)
					}

					if len(col[i].Title) > m.columnWidth {
						view += headlineStyle.Render(col[i].Title[:m.columnWidth-5] + "...")
					} else {
						view += headlineStyle.Render(col[i].Title)
					}
				} else {
					view += strings.Repeat(" ", m.columnWidth)
				}
				if j < len(m.client.Page.HeadlineColumns)-1 {
					view += cs
				}
			}
			view += "\n"
		}
	}
	return view
}
