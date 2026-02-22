package main

import tea "github.com/charmbracelet/bubbletea"

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		//The "left" and "h" keys move the cursor left
		case "left", "h":
			if m.cursorx > 0 {
				if (len(m.headlines[m.cursorx-1]) - 1) >= m.cursory {
					m.cursorx--
				}
			}

		// The "right" and "l" keys move the cursor right
		case "right", "l":
			if m.cursorx < len(m.headlines)-1 {
				if (len(m.headlines[m.cursorx+1]) - 1) >= m.cursory {
					m.cursorx++
				}
			}

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursory > 0 {
				m.cursory--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursory < m.curMaxRow-1 {
				m.cursory++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			//TODO: load page?
		}

	}

	m.curMaxRow = len(m.headlines[m.cursorx])
	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}
