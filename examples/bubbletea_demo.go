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

	// Header using character-based height
	headerNode := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   layout.Vw(100),
			Height:  layout.Ch(5), // 5 character heights
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

	// Content area with gradient
	contentHeight := m.height - 7 // Leave room for header + footer + margins
	if contentHeight > 0 {
		contentNode := &layout.Node{
			Style: layout.Style{
				Display:       layout.DisplayFlex,
				FlexDirection: layout.FlexDirectionRow,
				Width:         layout.Vw(100),
				Height:        layout.Px(float64(contentHeight)),
				Margin:        layout.Spacing{Top: layout.Ch(1), Right: layout.Px(0), Bottom: layout.Px(0), Left: layout.Px(0)},
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
					Width:   layout.Ch(1), // 1 character width
					Height:  layout.Px(float64(contentHeight)),
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
			Width:   layout.Vw(100),
			Height:  layout.Ch(1),
			Margin:  layout.Spacing{Top: layout.Ch(1), Right: layout.Px(0), Bottom: layout.Px(0), Left: layout.Px(0)},
		},
	}
	footerStyled := renderer.NewStyledNode(footerNode, nil)
	footerStyled.Content = "Press 'q' or ESC to quit • Resize terminal to see responsive layout • Now with CSS units!"
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
