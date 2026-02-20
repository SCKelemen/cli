package main

import (
	"fmt"
	"os"

	"github.com/SCKelemen/cli/renderer"
	"github.com/SCKelemen/color"
	"github.com/SCKelemen/layout"
	tea "github.com/charmbracelet/bubbletea"
)

type simpleModel struct {
	width  int
	height int
	ready  bool
}

func initialSimpleModel() simpleModel {
	return simpleModel{}
}

func (m simpleModel) Init() tea.Cmd {
	return nil
}

func (m simpleModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m simpleModel) View() string {
	if !m.ready {
		return "Initializing..."
	}

	screen := renderer.NewScreen(m.width, m.height)
	ctx := layout.NewLayoutContext(float64(m.width), float64(m.height), 16)

	// Root container
	root := &layout.Node{
		Style: layout.Style{
			Display:       layout.DisplayFlex,
			FlexDirection: layout.FlexDirectionColumn,
			Width:         layout.Vw(100),
			Height:        layout.Vh(100),
		},
	}
	rootStyled := renderer.NewStyledNode(root, nil)

	// Header - 15% of viewport height with minimum constraint
	headerNode := &layout.Node{
		Style: layout.Style{
			Display:   layout.DisplayBlock,
			Width:     layout.Vw(100),
			Height:    layout.Vh(15),
			MinHeight: layout.Ch(3), // Ensure minimum 3 rows
		},
	}
	purple, _ := color.ParseColor("#7D56F4")
	white, _ := color.ParseColor("#FAFAFA")
	headerStyle := &renderer.Style{
		Foreground:   &white,
		BorderColor:  &purple,
		TextOverflow: renderer.TextOverflowEllipsis, // Truncate long text with …
	}
	headerStyle.WithBorder(renderer.RoundedBorder)
	headerStyled := renderer.NewStyledNode(headerNode, headerStyle)
	headerStyled.Content = fmt.Sprintf("\n Responsive TUI Demo with CSS Units and Text Overflow!\n Terminal Size: %dx%d columns x rows", m.width, m.height)
	rootStyled.AddChild(headerStyled)

	// Content area - 75% of viewport height with gradient background
	contentNode := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   layout.Vw(100),
			Height:  layout.Vh(75),
		},
	}

	// Create gradient background by calculating color for this position
	gradientHue := float64(m.width*m.height) * 0.1
	for gradientHue >= 360 {
		gradientHue -= 360
	}
	gradientColor, _ := color.ParseColor(fmt.Sprintf("oklch(0.65 0.2 %.0f)", gradientHue))

	contentStyle := &renderer.Style{
		Foreground:  &white,
		Background:  &gradientColor,
		BorderColor: &gradientColor,
	}
	contentStyle.WithBorder(renderer.RoundedBorder)
	contentStyled := renderer.NewStyledNode(contentNode, contentStyle)
	contentStyled.Content = fmt.Sprintf(`

   Layout Engine with CSS Units!

   Viewport Units:
   • Width: Vw(100) = %d columns
   • Height: Vh(100) = %d rows
   • This content area: Vh(75)

   Text Units:
   • Ch(n) = character widths
   • Rem(n) = root font size

   Try resizing your terminal!
`, m.width, m.height)
	rootStyled.AddChild(contentStyled)

	// Footer - 10% of viewport height
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
	footerStyled.Content = "\n Press 'q' or ESC to quit • Resize to see responsive layout"
	rootStyled.AddChild(footerStyled)

	// Layout and render
	constraints := layout.Tight(float64(m.width), float64(m.height))
	layout.Layout(root, constraints, ctx)
	screen.Render(rootStyled)

	return screen.String()
}

func main() {
	p := tea.NewProgram(initialSimpleModel(), tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
