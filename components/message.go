package components

import (
	"github.com/SCKelemen/color"
	"github.com/SCKelemen/layout"
	"github.com/SCKelemen/cli/renderer"
)

// MessageBlock represents a styled message box
type MessageBlock struct {
	Text        string
	Foreground  *color.Color
	Background  *color.Color
	Border      renderer.BorderChars
	BorderColor *color.Color
}

// NewMessageBlock creates a new message block with default styling
func NewMessageBlock(text string) *MessageBlock {
	fg, _ := color.ParseColor("#FAFAFA")
	bc, _ := color.ParseColor("#7D56F4")
	return &MessageBlock{
		Text:        text,
		Foreground:  &fg,
		Background:  nil,
		Border:      renderer.RoundedBorder,
		BorderColor: &bc,
	}
}

// WithForeground sets the foreground color
func (m *MessageBlock) WithForeground(c *color.Color) *MessageBlock {
	m.Foreground = c
	return m
}

// WithBackground sets the background color
func (m *MessageBlock) WithBackground(c *color.Color) *MessageBlock {
	m.Background = c
	return m
}

// WithBorder sets the border style
func (m *MessageBlock) WithBorder(border renderer.BorderChars) *MessageBlock {
	m.Border = border
	return m
}

// WithBorderColor sets the border color
func (m *MessageBlock) WithBorderColor(c *color.Color) *MessageBlock {
	m.BorderColor = c
	return m
}

// ToStyledNode converts the message block to a styled node
func (m *MessageBlock) ToStyledNode() *renderer.StyledNode {
	// Calculate content dimensions (simple estimation)
	lines := 1
	maxWidth := len(m.Text)
	for _, c := range m.Text {
		if c == '\n' {
			lines++
		}
	}

	// Add padding for border
	width := float64(maxWidth + 4) // 2 for border + 2 for padding
	height := float64(lines + 2)   // 2 for border

	node := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   width,
			Height:  height,
		},
	}

	style := &renderer.Style{
		Foreground:  m.Foreground,
		Background:  m.Background,
		BorderColor: m.BorderColor,
	}
	style.WithBorder(m.Border)

	styledNode := renderer.NewStyledNode(node, style)
	styledNode.Content = m.Text

	return styledNode
}
