package main

import (
	"fmt"
	"os"

	"github.com/SCKelemen/cli/renderer"
	"github.com/SCKelemen/color"
	"github.com/SCKelemen/layout"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	minWidthTwoColumn   = 80  // Minimum width for 2-column layout
	minWidthThreeColumn = 120 // Minimum width for 3-column layout
	minPanelWidth       = 30  // Minimum width for each panel
)

type responsiveModel struct {
	width  int
	height int
	ready  bool
}

func initialResponsiveModel() responsiveModel {
	return responsiveModel{}
}

func (m responsiveModel) Init() tea.Cmd {
	return nil
}

func (m responsiveModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m responsiveModel) View() string {
	if !m.ready {
		return "Initializing..."
	}

	screen := renderer.NewScreen(m.width, m.height)

	// Determine layout mode based on width
	var layoutMode string
	var columns int
	if m.width >= minWidthThreeColumn {
		layoutMode = "Three Column"
		columns = 3
	} else if m.width >= minWidthTwoColumn {
		layoutMode = "Two Column"
		columns = 2
	} else {
		layoutMode = "Single Column"
		columns = 1
	}

	// Create root container
	root := &layout.Node{
		Style: layout.Style{
			Display:       layout.DisplayFlex,
			FlexDirection: layout.FlexDirectionColumn,
			Width:         layout.Px(float64(m.width)),
			Height:        layout.Px(float64(m.height)),
		},
	}
	rootStyled := renderer.NewStyledNode(root, nil)

	// Header showing current layout mode
	headerNode := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   layout.Px(float64(m.width)),
			Height:  layout.Px(3),
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
	headerStyled.Content = fmt.Sprintf(" Responsive Layout: %s • %dx%d", layoutMode, m.width, m.height)
	rootStyled.AddChild(headerStyled)

	// Content container with flexible layout
	contentHeight := m.height - 6 // Header + footer + margins
	if contentHeight > 0 {
		contentNode := &layout.Node{
			Style: layout.Style{
				Display:       layout.DisplayFlex,
				FlexDirection: layout.FlexDirectionRow,
				FlexWrap:      layout.FlexWrapWrap,
				Width:         layout.Px(float64(m.width)),
				Height:        layout.Px(float64(contentHeight)),
				Margin:        layout.Spacing{Top: layout.Px(1), Right: layout.Px(0), Bottom: layout.Px(1), Left: layout.Px(0)},
			},
		}
		contentStyled := renderer.NewStyledNode(contentNode, nil)

		// Create 6 panels that will wrap based on available width
		panelColors := []string{
			"#FF6B6B", // Red
			"#4ECDC4", // Cyan
			"#45B7D1", // Blue
			"#FFA07A", // Orange
			"#98D8C8", // Teal
			"#F7DC6F", // Yellow
		}

		panelTitles := []string{
			"Panel 1",
			"Panel 2",
			"Panel 3",
			"Panel 4",
			"Panel 5",
			"Panel 6",
		}

		// Calculate panel width based on columns
		panelWidth := m.width / columns
		if columns > 1 {
			panelWidth-- // Account for spacing
		}

		for i := 0; i < 6; i++ {
			panelNode := &layout.Node{
				Style: layout.Style{
					Display: layout.DisplayBlock,
					Width:   layout.Px(float64(panelWidth)),
					Height:  layout.Px(float64(contentHeight / 2)), // Two rows
					Margin:  layout.Spacing{Top: layout.Px(0), Right: layout.Px(1), Bottom: layout.Px(1), Left: layout.Px(0)},
				},
			}

			borderColor, _ := color.ParseColor(panelColors[i])
			panelStyle := &renderer.Style{
				Foreground:  &white,
				BorderColor: &borderColor,
			}
			panelStyle.WithBorder(renderer.RoundedBorder)
			panelStyled := renderer.NewStyledNode(panelNode, panelStyle)
			panelStyled.Content = fmt.Sprintf("\n %s\n\n Width: %d\n Mode: %s",
				panelTitles[i], panelWidth, layoutMode)
			contentStyled.AddChild(panelStyled)
		}

		rootStyled.AddChild(contentStyled)
	}

	// Footer with breakpoint info
	footerNode := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   layout.Px(float64(m.width)),
			Height:  layout.Px(2),
		},
	}
	gray, _ := color.ParseColor("#888888")
	footerStyle := &renderer.Style{
		Foreground: &gray,
	}
	footerStyled := renderer.NewStyledNode(footerNode, footerStyle)

	var breakpointInfo string
	if m.width >= minWidthThreeColumn {
		breakpointInfo = fmt.Sprintf("Desktop Mode (>=%d cols)", minWidthThreeColumn)
	} else if m.width >= minWidthTwoColumn {
		breakpointInfo = fmt.Sprintf("Tablet Mode (%d-%d cols)", minWidthTwoColumn, minWidthThreeColumn-1)
	} else {
		breakpointInfo = fmt.Sprintf("Mobile Mode (<%d cols)", minWidthTwoColumn)
	}

	footerStyled.Content = fmt.Sprintf("Breakpoint: %s • Resize terminal to see responsive behavior • q/ESC to quit", breakpointInfo)
	rootStyled.AddChild(footerStyled)

	// Layout and render
	constraints := layout.Tight(float64(m.width), float64(m.height))
	ctx := &layout.LayoutContext{
		ViewportWidth:  float64(m.width),
		ViewportHeight: float64(m.height),
		RootFontSize:   16,
	}
	layout.Layout(root, constraints, ctx)
	screen.Render(rootStyled)

	return screen.String()
}

func main() {
	p := tea.NewProgram(initialResponsiveModel(), tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
