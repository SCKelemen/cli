package main

import (
	"fmt"
	"os"

	"github.com/SCKelemen/cli/renderer"
	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/term"
)

type dimensionsModel struct {
	width  int
	height int
	dims   renderer.TerminalDimensions
	ready  bool
}

func initialDimensionsModel(dims renderer.TerminalDimensions) dimensionsModel {
	return dimensionsModel{
		dims: dims,
	}
}

func (m dimensionsModel) Init() tea.Cmd {
	return nil
}

func (m dimensionsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Update character grid dimensions
		m.dims.Columns = m.width
		m.dims.Rows = m.height

		// Recalculate pixel dimensions based on cell size
		// (don't re-query to avoid interfering with Bubbletea's terminal state)
		m.dims.PixelWidth = int(m.dims.CellWidth * float64(m.width))
		m.dims.PixelHeight = int(m.dims.CellHeight * float64(m.height))

		m.ready = true
		return m, nil
	}

	return m, nil
}

func (m dimensionsModel) View() string {
	if !m.ready {
		return "Initializing...\n"
	}

	var view string
	view += "Terminal Dimensions Test\n"
	view += "========================\n\n"

	view += fmt.Sprintf("Character Grid: %d columns × %d rows\n", m.dims.Columns, m.dims.Rows)
	view += "\n"

	if m.dims.HasPixelSupport {
		view += fmt.Sprintf("Pixel Dimensions: %d × %d pixels\n", m.dims.PixelWidth, m.dims.PixelHeight)
		view += fmt.Sprintf("Cell Size: %.2f × %.2f pixels\n", m.dims.CellWidth, m.dims.CellHeight)
		view += "\n"
		view += "✓ Terminal supports pixel size queries\n"
	} else {
		view += fmt.Sprintf("Pixel Dimensions: %d × %d pixels (estimated)\n", m.dims.PixelWidth, m.dims.PixelHeight)
		view += fmt.Sprintf("Cell Size: %.2f × %.2f pixels (estimated)\n", m.dims.CellWidth, m.dims.CellHeight)
		view += "\n"
		view += "✗ Terminal does not support pixel size queries\n"
		view += "  Using typical monospace font dimensions\n"
	}

	view += "\n"
	view += "Summary: " + m.dims.String() + "\n"
	view += "\n"
	view += "Press 'q' to quit\n"
	view += "Try resizing the terminal window!\n"

	return view
}

func main() {
	// Query dimensions BEFORE starting Bubbletea
	// This avoids interfering with Bubbletea's terminal state management
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width, height = 80, 24 // fallback
	}

	dims := renderer.QueryTerminalDimensions(width, height)

	fmt.Printf("Initial query: %s\n", dims.String())
	fmt.Printf("Press Enter to start TUI...")
	fmt.Scanln()

	p := tea.NewProgram(initialDimensionsModel(dims), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
