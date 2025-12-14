package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/SCKelemen/color"
	"github.com/SCKelemen/layout"
	"github.com/SCKelemen/cli/renderer"
)

type flexwrapModel struct {
	width  int
	height int
	ready  bool
}

func initialFlexwrapModel() flexwrapModel {
	return flexwrapModel{}
}

func (m flexwrapModel) Init() tea.Cmd {
	return nil
}

func (m flexwrapModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m flexwrapModel) View() string {
	if !m.ready {
		return "Initializing..."
	}

	screen := renderer.NewScreen(m.width, m.height)

	// Create root container
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
			Height:  3,
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
	headerStyled.Content = fmt.Sprintf(" Flexbox Auto-Wrap Demo • %dx%d", m.width, m.height)
	rootStyled.AddChild(headerStyled)

	// Flexbox container with wrapping
	contentHeight := m.height - 5
	if contentHeight > 0 {
		contentNode := &layout.Node{
			Style: layout.Style{
				Display:       layout.DisplayFlex,
				FlexDirection: layout.FlexDirectionRow,
				FlexWrap:      layout.FlexWrapWrap, // Enable wrapping
				Width:         float64(m.width),
				Height:        float64(contentHeight),
				Margin:        layout.Spacing{Top: 1, Right: 0, Bottom: 0, Left: 0},
			},
		}
		contentStyled := renderer.NewStyledNode(contentNode, nil)

		// Create cards with fixed minimum width
		// These will automatically wrap to new rows as terminal narrows
		cards := []struct {
			title string
			color string
			info  string
		}{
			{"Card A", "#FF6B6B", "Fixed width\n25 columns"},
			{"Card B", "#4ECDC4", "Will wrap\nautomatically"},
			{"Card C", "#45B7D1", "Based on\navailable space"},
			{"Card D", "#FFA07A", "Try resizing\nyour terminal"},
			{"Card E", "#98D8C8", "Flexbox\nwrapping"},
			{"Card F", "#F7DC6F", "Just like\nCSS flex-wrap"},
			{"Card G", "#E74C3C", "No media\nqueries needed"},
			{"Card H", "#3498DB", "Pure layout\nengine magic"},
		}

		cardWidth := 25   // Minimum card width
		cardHeight := 6   // Fixed card height

		for _, card := range cards {
			cardNode := &layout.Node{
				Style: layout.Style{
					Display:  layout.DisplayBlock,
					Width:    float64(cardWidth),
					Height:   float64(cardHeight),
					MinWidth: float64(cardWidth),
					Margin:   layout.Spacing{Top: 0, Right: 1, Bottom: 1, Left: 0},
				},
			}

			cardColor, _ := color.ParseColor(card.color)
			cardStyle := &renderer.Style{
				Foreground:  &white,
				BorderColor: &cardColor,
			}
			cardStyle.WithBorder(renderer.RoundedBorder)
			cardStyled := renderer.NewStyledNode(cardNode, cardStyle)
			cardStyled.Content = fmt.Sprintf("\n %s\n\n %s", card.title, card.info)
			contentStyled.AddChild(cardStyled)
		}

		rootStyled.AddChild(contentStyled)
	}

	// Footer
	footerNode := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   float64(m.width),
			Height:  1,
			Margin:  layout.Spacing{Top: 1, Right: 0, Bottom: 0, Left: 0},
		},
	}
	gray, _ := color.ParseColor("#888888")
	footerStyle := &renderer.Style{
		Foreground: &gray,
	}
	footerStyled := renderer.NewStyledNode(footerNode, footerStyle)

	cardsPerRow := m.width / 26 // cardWidth + margin
	if cardsPerRow < 1 {
		cardsPerRow = 1
	}
	footerStyled.Content = fmt.Sprintf("Cards per row: %d • Make terminal narrower to see wrapping • q/ESC to quit", cardsPerRow)
	rootStyled.AddChild(footerStyled)

	// Layout and render
	constraints := layout.Tight(float64(m.width), float64(m.height))
	layout.Layout(root, constraints)
	screen.Render(rootStyled)

	return screen.String()
}

func main() {
	p := tea.NewProgram(initialFlexwrapModel(), tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
