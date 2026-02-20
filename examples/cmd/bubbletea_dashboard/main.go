package main

import (
	"fmt"
	"os"
	"time"

	"github.com/SCKelemen/cli/renderer"
	"github.com/SCKelemen/color"
	"github.com/SCKelemen/layout"
	tea "github.com/charmbracelet/bubbletea"
)

type tickMsg time.Time

type dashboardModel struct {
	width   int
	height  int
	ready   bool
	counter int
}

func initialDashboardModel() dashboardModel {
	return dashboardModel{}
}

func (m dashboardModel) Init() tea.Cmd {
	return tickCmd()
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m dashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		return m, nil

	case tickMsg:
		m.counter++
		return m, tickCmd()
	}

	return m, nil
}

func (m dashboardModel) View() string {
	if !m.ready {
		return "Initializing..."
	}

	screen := renderer.NewScreen(m.width, m.height)

	// Create layout context
	ctx := layout.NewLayoutContext(float64(m.width), float64(m.height), 16)

	// Create responsive grid layout with viewport units
	root := &layout.Node{
		Style: layout.Style{
			Display:             layout.DisplayGrid,
			Width:               layout.Vw(100),
			Height:              layout.Vh(100),
			GridTemplateColumns: []layout.GridTrack{layout.FractionTrack(1), layout.FractionTrack(1)},
			GridTemplateRows:    []layout.GridTrack{layout.FixedTrack(layout.Vh(15)), layout.FractionTrack(1), layout.FractionTrack(1)},
			GridGap:             layout.Ch(1), // 1 character spacing between grid cells
		},
	}
	rootStyled := renderer.NewStyledNode(root, nil)

	// Header spanning both columns with character-based padding
	headerNode := &layout.Node{
		Style: layout.Style{
			Display:         layout.DisplayBlock,
			GridColumnStart: 1,
			GridColumnEnd:   3, // Span columns 1-2
			Padding:         layout.Uniform(layout.Ch(0.5)),
		},
	}
	purple, _ := color.ParseColor("#7D56F4")
	white, _ := color.ParseColor("#FAFAFA")
	headerStyle := &renderer.Style{
		Foreground:  &white,
		BorderColor: &purple,
	}
	headerStyle.WithBorder(renderer.RoundedBorder)
	headerStyled := renderer.NewStyledNode(headerNode, headerStyle)
	headerStyled.Content = fmt.Sprintf(" Dashboard (CSS Units!) • %dx%d • %ds", m.width, m.height, m.counter)
	rootStyled.AddChild(headerStyled)

	// Left panel - Color gradient
	leftNode := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
		},
	}
	cyan, _ := color.ParseColor("#00D7FF")
	leftStyle := &renderer.Style{
		Foreground:  &white,
		BorderColor: &cyan,
	}
	leftStyle.WithBorder(renderer.RoundedBorder)
	leftStyled := renderer.NewStyledNode(leftNode, leftStyle)
	leftStyled.Content = fmt.Sprintf("\n Gradient Panel\n\n Hue: %d°", (m.counter*10)%360)
	rootStyled.AddChild(leftStyled)

	// Right panel - Stats
	rightNode := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
		},
	}
	green, _ := color.ParseColor("#00FF87")
	rightStyle := &renderer.Style{
		Foreground:  &white,
		BorderColor: &green,
	}
	rightStyle.WithBorder(renderer.RoundedBorder)
	rightStyled := renderer.NewStyledNode(rightNode, rightStyle)
	rightStyled.Content = fmt.Sprintf("\n Statistics\n\n Width:  %d\n Height: %d\n Cells:  %d", m.width, m.height, m.width*m.height)
	rootStyled.AddChild(rightStyled)

	// Bottom left - Progress
	bottomLeftNode := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
		},
	}
	yellow, _ := color.ParseColor("#FFD700")
	bottomLeftStyle := &renderer.Style{
		Foreground:  &white,
		BorderColor: &yellow,
	}
	bottomLeftStyle.WithBorder(renderer.RoundedBorder)
	bottomLeftStyled := renderer.NewStyledNode(bottomLeftNode, bottomLeftStyle)
	progress := (m.counter % 10) * 10
	bottomLeftStyled.Content = fmt.Sprintf("\n Progress: %d%%", progress)
	rootStyled.AddChild(bottomLeftStyled)

	// Bottom right - Controls
	bottomRightNode := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
		},
	}
	red, _ := color.ParseColor("#FF5555")
	bottomRightStyle := &renderer.Style{
		Foreground:  &white,
		BorderColor: &red,
	}
	bottomRightStyle.WithBorder(renderer.RoundedBorder)
	bottomRightStyled := renderer.NewStyledNode(bottomRightNode, bottomRightStyle)
	bottomRightStyled.Content = "\n Controls\n\n q/ESC - Quit\n Resize terminal"
	rootStyled.AddChild(bottomRightStyled)

	// Layout and render with context
	constraints := layout.Tight(float64(m.width), float64(m.height))
	layout.Layout(root, constraints, ctx)
	screen.Render(rootStyled)

	return screen.String()
}

func main() {
	p := tea.NewProgram(initialDashboardModel(), tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
