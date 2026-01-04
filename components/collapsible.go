package components

import (
	"github.com/SCKelemen/color"
	"github.com/SCKelemen/layout"
	"github.com/SCKelemen/cli/renderer"
)

// Collapsible represents a collapsible section with a header and content
type Collapsible struct {
	Title       string
	Content     string
	Expanded    bool
	Foreground  *color.Color
	TitleColor  *color.Color
	Border      renderer.BorderChars
	BorderColor *color.Color
}

// NewCollapsible creates a new collapsible section
func NewCollapsible(title, content string) *Collapsible {
	fg, _ := color.ParseColor("#FAFAFA")
	tc, _ := color.ParseColor("#7D56F4")
	bc, _ := color.ParseColor("#5A5A5A")
	return &Collapsible{
		Title:       title,
		Content:     content,
		Expanded:    true,
		Foreground:  &fg,
		TitleColor:  &tc,
		Border:      renderer.RoundedBorder,
		BorderColor: &bc,
	}
}

// Toggle toggles the expanded state
func (c *Collapsible) Toggle() {
	c.Expanded = !c.Expanded
}

// ToStyledNode converts the collapsible to a styled node
func (c *Collapsible) ToStyledNode() *renderer.StyledNode {
	// Create root container
	root := &layout.Node{
		Style: layout.Style{
			Display:       layout.DisplayFlex,
			FlexDirection: layout.FlexDirectionColumn,
		},
	}

	rootStyled := renderer.NewStyledNode(root, nil)

	// Create header
	arrow := "▼"
	if !c.Expanded {
		arrow = "▶"
	}
	headerText := arrow + " " + c.Title

	headerNode := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   layout.Px(float64(len(headerText) + 4)),
			Height:  layout.Px(3), // Border + padding
			Padding: layout.Spacing{Top: layout.Px(0), Right: layout.Px(1), Bottom: layout.Px(0), Left: layout.Px(1)},
		},
	}

	headerStyle := &renderer.Style{
		Foreground:  c.TitleColor,
		Bold:        true,
		BorderColor: c.BorderColor,
	}
	headerStyle.WithBorder(c.Border)

	headerStyled := renderer.NewStyledNode(headerNode, headerStyle)
	headerStyled.Content = headerText
	rootStyled.AddChild(headerStyled)

	// Add content if expanded
	if c.Expanded && c.Content != "" {
		// Count lines in content
		lines := 1
		maxWidth := 0
		currentWidth := 0
		for _, ch := range c.Content {
			if ch == '\n' {
				lines++
				if currentWidth > maxWidth {
					maxWidth = currentWidth
				}
				currentWidth = 0
			} else {
				currentWidth++
			}
		}
		if currentWidth > maxWidth {
			maxWidth = currentWidth
		}

		contentNode := &layout.Node{
			Style: layout.Style{
				Display: layout.DisplayBlock,
				Width:   layout.Px(float64(maxWidth + 4)), // Add padding
				Height:  layout.Px(float64(lines + 2)),    // Add border
				Margin:  layout.Spacing{Top: layout.Px(1), Right: layout.Px(0), Bottom: layout.Px(0), Left: layout.Px(0)},
				Padding: layout.Spacing{Top: layout.Px(0), Right: layout.Px(1), Bottom: layout.Px(0), Left: layout.Px(1)},
			},
		}

		contentStyle := &renderer.Style{
			Foreground:  c.Foreground,
			BorderColor: c.BorderColor,
		}
		contentStyle.WithBorder(c.Border)

		contentStyled := renderer.NewStyledNode(contentNode, contentStyle)
		contentStyled.Content = c.Content
		rootStyled.AddChild(contentStyled)
	}

	return rootStyled
}
