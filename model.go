package main

import (
	_ "embed"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mikegetz/godrudge"
)

//go:embed logo
var logo string

const ansiReset = "\033[0m"
const ansiBold = "\033[1m"

type model struct {
	headlines [][]godrudge.Headline // all headlines of drudge report
	cursorx   int                   // the current column
	cursory   int                   // the current row in the current column
	curMaxRow int                   // the max number of rows in the current column
	width     int                   // width of the terminal
	maxRows   int                   // represents the column with the most headlines
	selected  godrudge.Headline     // the currently selected headline
}

func initialModel() model {
	client := godrudge.NewClient()
	err := client.ParseRSS()
	if err != nil {
		fmt.Printf("Error parsing drudgereport: %v\n", err)
		os.Exit(1)
	}
	model := model{
		headlines: client.Page.HeadlineColumns,
		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
	}

	maxRows := 0
	for _, col := range model.headlines {
		if len(col) > maxRows {
			maxRows = len(col)
		}
	}

	model.maxRows = maxRows

	return model
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}
