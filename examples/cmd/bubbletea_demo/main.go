package main

import (
	"fmt"
	"os"

	"github.com/SCKelemen/cli/renderer"
	"github.com/SCKelemen/color"
	"github.com/SCKelemen/layout"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	width  int
	height int
	ready  bool
}

func initialModel() model {
	return model{}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
	}

	return m, nil
}

func (m model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	// Create screen buffer matching terminal size
	screen := renderer.NewScreen(m.width, m.height)

	// Create layout context with viewport dimensions and root font size
	ctx := layout.NewLayoutContext(float64(m.width), float64(m.height), 16)

	// Create root layout node using viewport units
	root := &layout.Node{
		Style: layout.Style{
			Display:       layout.DisplayFlex,
			FlexDirection: layout.FlexDirectionColumn,
			Width:         layout.Vw(100), // 100% of viewport width
			Height:        layout.Vh(100), // 100% of viewport height
		},
	}
	rootStyled := renderer.NewStyledNode(root, nil)

	// Header - 15% of viewport height with minimum
	headerNode := &layout.Node{
		Style: layout.Style{
			Display:   layout.DisplayBlock,
			Width:     layout.Vw(100),
			Height:    layout.Vh(15),
			MinHeight: layout.Ch(3), // Minimum 3 character heights
		},
	}
	headerColor, _ := color.ParseColor("#7D56F4")
	textColor, _ := color.ParseColor("#FAFAFA")
	headerStyle := &renderer.Style{
		Foreground:  &textColor,
		BorderColor: &headerColor,
	}
	headerStyle.WithBorder(renderer.RoundedBorder)
	headerStyled := renderer.NewStyledNode(headerNode, headerStyle)
	headerStyled.Content = fmt.Sprintf("\n Responsive TUI • %dx%d", m.width, m.height)
	rootStyled.AddChild(headerStyled)

	// Content area - 75% of viewport with gradient background
	// Create a single panel with gradient background color
	contentNode := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   layout.Vw(100),
			Height:  layout.Vh(75),
		},
	}

	// Rainbow gradient - calculate hue based on time/position
	hue := float64((m.width + m.height) % 360)
	gradientColor, _ := color.ParseColor(fmt.Sprintf("oklch(0.65 0.2 %.0f)", hue))

	contentStyle := &renderer.Style{
		Foreground:  &textColor,
		Background:  &gradientColor,
		BorderColor: &gradientColor,
	}
	contentStyle.WithBorder(renderer.RoundedBorder)
	contentStyled := renderer.NewStyledNode(contentNode, contentStyle)

	// Simpler content that fits in small viewports
	contentStyled.Content = fmt.Sprintf("\n\n %dx%d • Vh(%d/%d/%d)\n Try resizing!",
		m.width, m.height, 15, 75, 10)
	rootStyled.AddChild(contentStyled)

	// Footer - 10% of viewport
	footerNode := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   layout.Vw(100),
			Height:  layout.Vh(10),
		},
	}
	gray, _ := color.ParseColor("#888888")
	footerStyle := &renderer.Style{
		Foreground: &gray,
	}
	footerStyled := renderer.NewStyledNode(footerNode, footerStyle)
	footerStyled.Content = " Press 'q' to quit • Resize window"
	rootStyled.AddChild(footerStyled)

	// Layout and render with context
	constraints := layout.Tight(float64(m.width), float64(m.height))
	layout.Layout(root, constraints, ctx)
	screen.Render(rootStyled)

	return screen.String()
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
