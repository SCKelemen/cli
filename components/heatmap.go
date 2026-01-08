package components

import (
	"github.com/SCKelemen/cli/renderer"
	"github.com/SCKelemen/color"
	"github.com/SCKelemen/dataviz"
	design "github.com/SCKelemen/design-system"
	"github.com/SCKelemen/layout"
)

// Heatmap represents a contribution heatmap visualization component
type Heatmap struct {
	Data         dataviz.HeatmapData
	Width        int
	Height       int
	Color        string
	Theme        string
	DesignTokens *design.DesignTokens
}

// NewHeatmap creates a new heatmap component with default settings
func NewHeatmap(data dataviz.HeatmapData) *Heatmap {
	return &Heatmap{
		Data:         data,
		Width:        80,
		Height:       7,
		Color:        "#3B82F6",
		Theme:        "default",
		DesignTokens: design.DefaultTheme(),
	}
}

// WithSize sets the width and height
func (h *Heatmap) WithSize(width, height int) *Heatmap {
	h.Width = width
	h.Height = height
	return h
}

// WithColor sets the primary color
func (h *Heatmap) WithColor(c string) *Heatmap {
	h.Color = c
	return h
}

// WithTheme sets the theme
func (h *Heatmap) WithTheme(theme string) *Heatmap {
	h.Theme = theme
	switch theme {
	case "midnight":
		h.DesignTokens = design.MidnightTheme()
	case "nord":
		h.DesignTokens = design.NordTheme()
	case "paper":
		h.DesignTokens = design.PaperTheme()
	case "wrapped":
		h.DesignTokens = design.WrappedTheme()
	default:
		h.DesignTokens = design.DefaultTheme()
	}
	return h
}

// WithDesignTokens sets custom design tokens
func (h *Heatmap) WithDesignTokens(tokens *design.DesignTokens) *Heatmap {
	h.DesignTokens = tokens
	return h
}

// ToStyledNode converts the heatmap to a styled node for terminal rendering
func (h *Heatmap) ToStyledNode() *renderer.StyledNode {
	// Create bounds for rendering
	bounds := dataviz.Bounds{
		X:      0,
		Y:      0,
		Width:  h.Width,
		Height: h.Height,
	}

	// Create render config
	config := dataviz.RenderConfig{
		DesignTokens: h.DesignTokens,
		Color:        h.Color,
		Theme:        h.Theme,
	}

	// Render using terminal renderer
	termRenderer := dataviz.NewTerminalRenderer()
	output := termRenderer.RenderHeatmap(h.Data, bounds, config)

	// Create layout node
	node := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   layout.Px(float64(h.Width)),
			Height:  layout.Px(float64(h.Height)),
		},
	}

	// Parse theme colors for styling
	fg, _ := color.ParseColor(h.DesignTokens.Color)
	bg, _ := color.ParseColor(h.DesignTokens.Background)

	style := &renderer.Style{
		Foreground: &fg,
		Background: &bg,
	}

	styledNode := renderer.NewStyledNode(node, style)
	styledNode.Content = output.String()

	return styledNode
}
