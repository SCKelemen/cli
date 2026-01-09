package components

import (
	"github.com/SCKelemen/cli/renderer"
	"github.com/SCKelemen/color"
	"github.com/SCKelemen/dataviz"
	design "github.com/SCKelemen/design-system"
	"github.com/SCKelemen/layout"
)

// AreaChart represents an area chart visualization component
type AreaChart struct {
	Data         dataviz.AreaChartData
	Width        int
	Height       int
	Color        string
	Theme        string
	DesignTokens *design.DesignTokens
}

// NewAreaChart creates a new area chart component with default settings
func NewAreaChart(data dataviz.AreaChartData) *AreaChart {
	return &AreaChart{
		Data:         data,
		Width:        80,
		Height:       20,
		Color:        "#3B82F6",
		Theme:        "default",
		DesignTokens: design.DefaultTheme(),
	}
}

// WithSize sets the width and height
func (a *AreaChart) WithSize(width, height int) *AreaChart {
	a.Width = width
	a.Height = height
	return a
}

// WithColor sets the primary color
func (a *AreaChart) WithColor(c string) *AreaChart {
	a.Color = c
	return a
}

// WithTheme sets the theme
func (a *AreaChart) WithTheme(theme string) *AreaChart {
	a.Theme = theme
	switch theme {
	case "midnight":
		a.DesignTokens = design.MidnightTheme()
	case "nord":
		a.DesignTokens = design.NordTheme()
	case "paper":
		a.DesignTokens = design.PaperTheme()
	case "wrapped":
		a.DesignTokens = design.WrappedTheme()
	default:
		a.DesignTokens = design.DefaultTheme()
	}
	return a
}

// WithDesignTokens sets custom design tokens
func (a *AreaChart) WithDesignTokens(tokens *design.DesignTokens) *AreaChart {
	a.DesignTokens = tokens
	return a
}

// ToStyledNode converts the area chart to a styled node for terminal rendering
func (a *AreaChart) ToStyledNode() *renderer.StyledNode {
	// Create bounds for rendering
	bounds := dataviz.Bounds{
		X:      0,
		Y:      0,
		Width:  a.Width,
		Height: a.Height,
	}

	// Create render config
	config := dataviz.RenderConfig{
		DesignTokens: a.DesignTokens,
		Color:        a.Color,
		Theme:        a.Theme,
	}

	// Render using terminal renderer
	termRenderer := dataviz.NewTerminalRenderer()
	output := termRenderer.RenderAreaChart(a.Data, bounds, config)

	// Create layout node
	node := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   layout.Px(float64(a.Width)),
			Height:  layout.Px(float64(a.Height)),
		},
	}

	// Parse theme colors for styling
	fg, _ := color.ParseColor(a.DesignTokens.Color)
	bg, _ := color.ParseColor(a.DesignTokens.Background)

	style := &renderer.Style{
		Foreground: &fg,
		Background: &bg,
	}

	styledNode := renderer.NewStyledNode(node, style)
	styledNode.Content = output.String()

	return styledNode
}

// Render converts the area chart to a string for terminal display
func (a *AreaChart) Render() string {
	node := a.ToStyledNode()
	screen := renderer.NewScreen(a.Width, a.Height)
	screen.Render(node)
	return screen.String()
}
