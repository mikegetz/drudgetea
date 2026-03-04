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
	cs               = "  "
	containerPadding = 2
	containerStyle   = lipgloss.NewStyle().Padding(0, containerPadding)
)

func (m model) View() tea.View {
	contentWidth := m.width - (containerPadding * 2)
	centerStyle := lipgloss.NewStyle().Width(contentWidth).Align(lipgloss.Center)

	view := "\n"

	view += m.TopHeadlineView(contentWidth)
	view += m.MainHeadlineView(contentWidth)

	if m.showLogo {
		view += centerStyle.Render(logo) + "\n"
	} else {
		view += "\n"
	}

	view += m.ColumnView(contentWidth)

	// The footer
	view += m.FooterView(contentWidth)

	// The help view
	m.help.SetWidth(contentWidth)
	helpView := m.help.View(m.keys)
	height := 4 - strings.Count(helpView, "\n")

	view += strings.Repeat("\n", height) + helpView + "\n"

	// Send the UI for rendering

	teaView := tea.NewView(containerStyle.Render(view))

	teaView.BackgroundColor = lipgloss.Color("#181A1B")

	return teaView
}

func (m model) FooterView(contentWidth int) string {
	footerStyleColor := m.cursorStyle.UnsetAlign().UnsetWidth()
	footerStyleBorder := lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true).BorderForeground(lipgloss.Color("#888888")).Width(contentWidth)

	// Debug info
	var debug string
	if os.Getenv("DEBUG") != "" {
		debug = m.inputStyle.Render("\nDebug: " + fmt.Sprintf("cursorx: %d, cursory: %d, curMaxRow: %d, columnWidth: %d, foregroundColor: %v, time: %s", m.cursorx, m.cursory, m.curMaxRow, m.columnWidth, m.selected.Style.GetForeground(), m.time.Format("15:04:05")))
	}

	view := footerStyleColor.Render(m.selected.Title)

	urlLine := "(" + lipgloss.NewStyle().Hyperlink(m.selected.URL).Render(m.selected.URL) + ")"
	view += "\n" + urlLine

	if m.toggleRowLess == 0 || os.Getenv("DEBUG") != "" {
		timeSince := time.Since(m.time).Truncate(time.Second).String()
		timeStr := m.help.Styles.ShortKey.Render(timeSince) + m.help.Styles.ShortDesc.Render(" since last refresh")
		view += "\n" + timeStr
	}

	view += debug

	return footerStyleBorder.Render(view)
}

func (m *model) TopHeadlineView(contentWidth int) string {
	var topHeadlines []godrudge.Headline
	if m.toggleRowLess > 0 && m.toggleRowLess < len(m.client.Page.TopHeadlines) {
		topHeadlines = m.client.Page.TopHeadlines[:m.toggleRowLess]
	} else {
		topHeadlines = m.client.Page.TopHeadlines
	}
	view := ""
	for i, topHeadline := range topHeadlines {
		headlineStyle := topHeadline.Style.Width(contentWidth).Align(lipgloss.Left).Hyperlink(topHeadline.URL)

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

func (m *model) MainHeadlineView(contentWidth int) string {
	var mainHeadlines []godrudge.Headline
	if m.toggleRowLess > 0 && m.toggleRowLess < len(m.client.Page.MainHeadlines) {
		mainHeadlines = m.client.Page.MainHeadlines[:m.toggleRowLess]
	} else {
		mainHeadlines = m.client.Page.MainHeadlines
	}
	view := ""
	for i, mainHeadline := range mainHeadlines {
		var headlineStyle = mainHeadline.Style.Width(contentWidth).Align(lipgloss.Center).Hyperlink(mainHeadline.URL)

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

func (m *model) ColumnView(contentWidth int) string {
	columnContentWidth := contentWidth - (containerPadding * 2)
	var rows int
	if m.toggleRowLess > 0 {
		rows = m.toggleRowLess
	} else {
		rows = m.maxRows
	}

	m.columnWidth = (columnContentWidth - (len(cs) * 2)) / len(m.client.Page.HeadlineColumns)
	view := ""
	if m.columnWidth > 3 {
		for i := 0; i < rows; i++ {
			for colIndex, col := range m.client.Page.HeadlineColumns {
				if i < len(col) {
					headlineStyle := col[i].Style.Width(m.columnWidth).Hyperlink(col[i].URL)

					if m.disableLinkgloss {
						headlineStyle = headlineStyle.UnsetHyperlink()
					}

					switch colIndex {
					case 0:
						headlineStyle = headlineStyle.Align(lipgloss.Left)
					case 1:
						headlineStyle = headlineStyle.Align(lipgloss.Center)
					case 2:
						headlineStyle = headlineStyle.Align(lipgloss.Right)
					}

					if i == m.cursory && colIndex == m.cursorx && m.cursorGroup == 2 {
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
				if colIndex < len(m.client.Page.HeadlineColumns)-1 {
					view += cs
				}
			}
			view += "\n"
		}
	}
	return containerStyle.Render(view) + "\n"
}
