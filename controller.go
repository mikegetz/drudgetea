package main

import (
	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width

	case tea.KeyMsg:
		switch {

		// These keys should exit the program.
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		//The "left" and "h" keys move the cursor left
		case key.Matches(msg, m.keys.Left):
			if m.cursorGroup == 2 {
				if m.cursorx > 0 {
					if (len(m.headlines[m.cursorx-1]) - 1) >= m.cursory {
						m.cursorx--
					}
				}
			}

		// The "right" and "l" keys move the cursor right
		case key.Matches(msg, m.keys.Right):
			if m.cursorGroup == 2 {
				if m.cursorx < len(m.headlines)-1 {
					if (len(m.headlines[m.cursorx+1]) - 1) >= m.cursory {
						m.cursorx++
					}
				}
			}

		// The "up" and "k" keys move the cursor up
		case key.Matches(msg, m.keys.Up):
			if m.cursory > 0 {
				m.cursory--
			}

		// The "down" and "j" keys move the cursor down
		case key.Matches(msg, m.keys.Down):
			if m.cursory < m.curMaxRow-1 {
				m.cursory++
			}

		// ? toggles the help view
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll

		// c copy to clipboard
		case key.Matches(msg, m.keys.Copy):
			if m.selected.Href != "" {
				return m, tea.Batch(copyToClipboardCmd(m.selected.Href))
			}

		// tab switches headline group
		case key.Matches(msg, m.keys.Tab):
			m.cursorGroup++
			if m.cursorGroup > 2 {
				m.cursorGroup = 0
			}
			m.cursory = 0

		// l toggles ascii art logo
		case key.Matches(msg, m.keys.ShowLogo):
			m.showLogo = !m.showLogo

		// space or enter toggles more or less rows
		case key.Matches(msg, m.keys.Less):
			if m.toggleRowLess == 0 {
				m.toggleRowLess = 10
			} else {
				m.toggleRowLess = 0
			}
		}

	}

	// Update the selected headline based on the new cursor position.
	// Track ColumnHeadline Row length for cursor movement and column expansion
	switch m.cursorGroup {
	case 0:
		m.selected = m.topHeadlines[m.cursory]
		m.curMaxRow = len(m.topHeadlines)
	case 1:
		m.selected = m.mainHeadlines[m.cursory]
		m.curMaxRow = len(m.mainHeadlines)
	case 2:
		m.selected = m.headlines[m.cursorx][m.cursory]
		if m.toggleRowLess == 0 {
			m.curMaxRow = len(m.headlines[m.cursorx])
		} else {
			m.curMaxRow = m.toggleRowLess
		}
	}
	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

type clipboardMsg struct {
	err error
}

func copyToClipboardCmd(s string) tea.Cmd {
	return func() tea.Msg {
		err := clipboard.WriteAll(s)
		return clipboardMsg{err: err}
	}
}
