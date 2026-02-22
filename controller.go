package main

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch {

		// These keys should exit the program.
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		//The "left" and "h" keys move the cursor left
		case key.Matches(msg, m.keys.Left):
			if m.cursorx > 0 {
				if (len(m.headlines[m.cursorx-1]) - 1) >= m.cursory {
					m.cursorx--
				}
			}

		// The "right" and "l" keys move the cursor right
		case key.Matches(msg, m.keys.Right):
			if m.cursorx < len(m.headlines)-1 {
				if (len(m.headlines[m.cursorx+1]) - 1) >= m.cursory {
					m.cursorx++
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

		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll

		// the selected state for the item that the cursor is pointing at.
		case key.Matches(msg, m.keys.Select):
			//TODO: load page?
		}

	}

	m.curMaxRow = len(m.headlines[m.cursorx])
	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}
