package components

import (
	"github.com/SCKelemen/cli/renderer"
	"github.com/SCKelemen/color"
	"github.com/SCKelemen/dataviz"
	design "github.com/SCKelemen/design-system"
	"github.com/SCKelemen/layout"
)

// LineGraph represents a line graph visualization component
type LineGraph struct {
	Data         dataviz.LineGraphData
	Width        int
	Height       int
	Color        string
	Theme        string
	DesignTokens *design.DesignTokens
}

// NewLineGraph creates a new line graph component with default settings
func NewLineGraph(data dataviz.LineGraphData) *LineGraph {
	return &LineGraph{
		Data:         data,
		Width:        80,
		Height:       20,
		Color:        "#3B82F6",
		Theme:        "default",
		DesignTokens: design.DefaultTheme(),
	}
}

// WithSize sets the width and height
func (l *LineGraph) WithSize(width, height int) *LineGraph {
	l.Width = width
	l.Height = height
	return l
}

// WithColor sets the primary color
func (l *LineGraph) WithColor(c string) *LineGraph {
	l.Color = c
	return l
}

// WithTheme sets the theme
func (l *LineGraph) WithTheme(theme string) *LineGraph {
	l.Theme = theme
	switch theme {
	case "midnight":
		l.DesignTokens = design.MidnightTheme()
	case "nord":
		l.DesignTokens = design.NordTheme()
	case "paper":
		l.DesignTokens = design.PaperTheme()
	case "wrapped":
		l.DesignTokens = design.WrappedTheme()
	default:
		l.DesignTokens = design.DefaultTheme()
	}
	return l
}

// WithDesignTokens sets custom design tokens
func (l *LineGraph) WithDesignTokens(tokens *design.DesignTokens) *LineGraph {
	l.DesignTokens = tokens
	return l
}

// ToStyledNode converts the line graph to a styled node for terminal rendering
func (l *LineGraph) ToStyledNode() *renderer.StyledNode {
	// Create bounds for rendering
	bounds := dataviz.Bounds{
		X:      0,
		Y:      0,
		Width:  l.Width,
		Height: l.Height,
	}

	// Create render config
	config := dataviz.RenderConfig{
		DesignTokens: l.DesignTokens,
		Color:        l.Color,
		Theme:        l.Theme,
	}

	// Render using terminal renderer
	termRenderer := dataviz.NewTerminalRenderer()
	output := termRenderer.RenderLineGraph(l.Data, bounds, config)

	// Create layout node
	node := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   layout.Px(float64(l.Width)),
			Height:  layout.Px(float64(l.Height)),
		},
	}

	// Parse theme colors for styling
	fg, _ := color.ParseColor(l.DesignTokens.Color)
	bg, _ := color.ParseColor(l.DesignTokens.Background)

	style := &renderer.Style{
		Foreground: &fg,
		Background: &bg,
	}

	styledNode := renderer.NewStyledNode(node, style)
	styledNode.Content = output.String()

	return styledNode
}
