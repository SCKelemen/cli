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

	// Create root layout node
	root := &layout.Node{
		Style: layout.Style{
			Display:       layout.DisplayFlex,
			FlexDirection: layout.FlexDirectionColumn,
			Width:         float64(m.width),
			Height:        float64(m.height),
		},
	}
	rootStyled := renderer.NewStyledNode(root, nil)

	// Header
	headerNode := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   float64(m.width),
			Height:  5,
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
	headerStyled.Content = fmt.Sprintf("\n Terminal UI Demo - Responsive Layout\n Size: %dx%d", m.width, m.height)
	rootStyled.AddChild(headerStyled)

	// Content area with gradient
	contentHeight := m.height - 7 // Leave room for header + footer + margins
	if contentHeight > 0 {
		contentNode := &layout.Node{
			Style: layout.Style{
				Display:       layout.DisplayFlex,
				FlexDirection: layout.FlexDirectionRow,
				Width:         float64(m.width),
				Height:        float64(contentHeight),
				Margin:        layout.Spacing{Top: 1, Right: 0, Bottom: 0, Left: 0},
			},
		}
		contentStyled := renderer.NewStyledNode(contentNode, nil)

		// Create gradient cells across full width
		for i := 0; i < m.width; i++ {
			t := float64(i) / float64(m.width-1)
			hue := t * 360
			colorStr := fmt.Sprintf("oklch(0.65 0.2 %.0f)", hue)
			c, _ := color.ParseColor(colorStr)

			cellNode := &layout.Node{
				Style: layout.Style{
					Display: layout.DisplayBlock,
					Width:   1,
					Height:  float64(contentHeight),
				},
			}
			cellStyle := &renderer.Style{
				Background: &c,
			}
			cellStyled := renderer.NewStyledNode(cellNode, cellStyle)
			cellStyled.Content = " "
			contentStyled.AddChild(cellStyled)
		}

		rootStyled.AddChild(contentStyled)
	}

	// Footer with instructions
	footerNode := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   float64(m.width),
			Height:  1,
			Margin:  layout.Spacing{Top: 1, Right: 0, Bottom: 0, Left: 0},
		},
	}
	footerStyled := renderer.NewStyledNode(footerNode, nil)
	footerStyled.Content = "Press 'q' or ESC to quit • Resize terminal to see responsive layout"
	rootStyled.AddChild(footerStyled)

	// Layout and render
	constraints := layout.Tight(float64(m.width), float64(m.height))
	layout.Layout(root, constraints)
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
