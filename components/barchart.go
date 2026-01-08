package components

import (
	"github.com/SCKelemen/cli/renderer"
	"github.com/SCKelemen/color"
	"github.com/SCKelemen/dataviz"
	design "github.com/SCKelemen/design-system"
	"github.com/SCKelemen/layout"
)

// BarChart represents a bar chart visualization component
type BarChart struct {
	Data         dataviz.BarChartData
	Width        int
	Height       int
	Color        string
	Theme        string
	DesignTokens *design.DesignTokens
}

// NewBarChart creates a new bar chart component with default settings
func NewBarChart(data dataviz.BarChartData) *BarChart {
	return &BarChart{
		Data:         data,
		Width:        60,
		Height:       len(data.Bars) + 2, // Auto-adjust height based on number of bars
		Color:        "#3B82F6",
		Theme:        "default",
		DesignTokens: design.DefaultTheme(),
	}
}

// WithSize sets the width and height
func (b *BarChart) WithSize(width, height int) *BarChart {
	b.Width = width
	b.Height = height
	return b
}

// WithColor sets the primary color
func (b *BarChart) WithColor(c string) *BarChart {
	b.Color = c
	return b
}

// WithTheme sets the theme
func (b *BarChart) WithTheme(theme string) *BarChart {
	b.Theme = theme
	switch theme {
	case "midnight":
		b.DesignTokens = design.MidnightTheme()
	case "nord":
		b.DesignTokens = design.NordTheme()
	case "paper":
		b.DesignTokens = design.PaperTheme()
	case "wrapped":
		b.DesignTokens = design.WrappedTheme()
	default:
		b.DesignTokens = design.DefaultTheme()
	}
	return b
}

// WithDesignTokens sets custom design tokens
func (b *BarChart) WithDesignTokens(tokens *design.DesignTokens) *BarChart {
	b.DesignTokens = tokens
	return b
}

// ToStyledNode converts the bar chart to a styled node for terminal rendering
func (b *BarChart) ToStyledNode() *renderer.StyledNode {
	// Create bounds for rendering
	bounds := dataviz.Bounds{
		X:      0,
		Y:      0,
		Width:  b.Width,
		Height: b.Height,
	}

	// Create render config
	config := dataviz.RenderConfig{
		DesignTokens: b.DesignTokens,
		Color:        b.Color,
		Theme:        b.Theme,
	}

	// Render using terminal renderer
	termRenderer := dataviz.NewTerminalRenderer()
	output := termRenderer.RenderBarChart(b.Data, bounds, config)

	// Create layout node
	node := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   layout.Px(float64(b.Width)),
			Height:  layout.Px(float64(b.Height)),
		},
	}

	// Parse theme colors for styling
	fg, _ := color.ParseColor(b.DesignTokens.Color)
	bg, _ := color.ParseColor(b.DesignTokens.Background)

	style := &renderer.Style{
		Foreground: &fg,
		Background: &bg,
	}

	styledNode := renderer.NewStyledNode(node, style)
	styledNode.Content = output.String()

	return styledNode
}
