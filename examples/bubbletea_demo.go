package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/SCKelemen/color"
	"github.com/SCKelemen/layout"
	"github.com/SCKelemen/cli/renderer"
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

	// Header - 15% of viewport height
	headerNode := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   layout.Vw(100),
			Height:  layout.Vh(15),
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
	headerStyled.Content = fmt.Sprintf("\n Terminal UI Demo - Responsive Layout\n Size: %dx%d • Using CSS Units!", m.width, m.height)
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
	contentStyled.Content = fmt.Sprintf(`


   Full Viewport Layout Demo

   Header: Vh(15) = %d rows
   Content: Vh(75) = %d rows
   Footer: Vh(10) = %d rows

   Total: 100%% of terminal height

   Try resizing to see proportions adapt!

`, int(float64(m.height)*0.15), int(float64(m.height)*0.75), int(float64(m.height)*0.10))
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
	footerStyled.Content = "\n Press 'q' or ESC to quit • Resize terminal to see responsive layout"
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
