package components

import (
	"github.com/SCKelemen/cli/renderer"
	"github.com/SCKelemen/color"
	"github.com/SCKelemen/dataviz"
	design "github.com/SCKelemen/design-system"
	"github.com/SCKelemen/layout"
)

// ScatterPlot represents a scatter plot visualization component
type ScatterPlot struct {
	Data         dataviz.ScatterPlotData
	Width        int
	Height       int
	Color        string
	Theme        string
	DesignTokens *design.DesignTokens
}

// NewScatterPlot creates a new scatter plot component with default settings
func NewScatterPlot(data dataviz.ScatterPlotData) *ScatterPlot {
	return &ScatterPlot{
		Data:         data,
		Width:        80,
		Height:       20,
		Color:        "#3B82F6",
		Theme:        "default",
		DesignTokens: design.DefaultTheme(),
	}
}

// WithSize sets the width and height
func (s *ScatterPlot) WithSize(width, height int) *ScatterPlot {
	s.Width = width
	s.Height = height
	return s
}

// WithColor sets the primary color
func (s *ScatterPlot) WithColor(c string) *ScatterPlot {
	s.Color = c
	return s
}

// WithTheme sets the theme
func (s *ScatterPlot) WithTheme(theme string) *ScatterPlot {
	s.Theme = theme
	switch theme {
	case "midnight":
		s.DesignTokens = design.MidnightTheme()
	case "nord":
		s.DesignTokens = design.NordTheme()
	case "paper":
		s.DesignTokens = design.PaperTheme()
	case "wrapped":
		s.DesignTokens = design.WrappedTheme()
	default:
		s.DesignTokens = design.DefaultTheme()
	}
	return s
}

// WithDesignTokens sets custom design tokens
func (s *ScatterPlot) WithDesignTokens(tokens *design.DesignTokens) *ScatterPlot {
	s.DesignTokens = tokens
	return s
}

// ToStyledNode converts the scatter plot to a styled node for terminal rendering
func (s *ScatterPlot) ToStyledNode() *renderer.StyledNode {
	// Create bounds for rendering
	bounds := dataviz.Bounds{
		X:      0,
		Y:      0,
		Width:  s.Width,
		Height: s.Height,
	}

	// Create render config
	config := dataviz.RenderConfig{
		DesignTokens: s.DesignTokens,
		Color:        s.Color,
		Theme:        s.Theme,
	}

	// Render using terminal renderer
	termRenderer := dataviz.NewTerminalRenderer()
	output := termRenderer.RenderScatterPlot(s.Data, bounds, config)

	// Create layout node
	node := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   layout.Px(float64(s.Width)),
			Height:  layout.Px(float64(s.Height)),
		},
	}

	// Parse theme colors for styling
	fg, _ := color.ParseColor(s.DesignTokens.Color)
	bg, _ := color.ParseColor(s.DesignTokens.Background)

	style := &renderer.Style{
		Foreground: &fg,
		Background: &bg,
	}

	styledNode := renderer.NewStyledNode(node, style)
	styledNode.Content = output.String()

	return styledNode
}

// Render converts the scatter plot to a string for terminal display
func (s *ScatterPlot) Render() string {
	node := s.ToStyledNode()
	screen := renderer.NewScreen(s.Width, s.Height)
	screen.Render(node)
	return screen.String()
}
