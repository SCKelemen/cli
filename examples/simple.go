package main

import (
	"fmt"

	"github.com/SCKelemen/color"
	"github.com/SCKelemen/layout"
	"github.com/SCKelemen/cli/components"
	"github.com/SCKelemen/cli/renderer"
)

func main() {
	// Terminal dimensions
	width, height := 80, 25

	// Create screen
	screen := renderer.NewScreen(width, height)

	// Build simple layout
	root := &layout.Node{
		Style: layout.Style{
			Display:       layout.DisplayFlex,
			FlexDirection: layout.FlexDirectionColumn,
			Width:         float64(width),
			Height:        float64(height),
			Padding:       layout.Uniform(2),
		},
	}
	rootStyled := renderer.NewStyledNode(root, nil)

	// Header
	headerNode := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   float64(width - 4),
			Height:  5,
		},
	}
	fgWhite, _ := color.ParseColor("#FFFFFF")
	bgPurple, _ := color.ParseColor("oklch(0.5 0.2 270)")
	borderPurple, _ := color.ParseColor("oklch(0.7 0.2 270)")
	headerStyle := &renderer.Style{
		Foreground:  &fgWhite,
		Background:  &bgPurple,
		Bold:        true,
		BorderColor: &borderPurple,
	}
	headerStyle.WithBorder(renderer.ThickBorder)
	headerStyled := renderer.NewStyledNode(headerNode, headerStyle)
	headerStyled.Content = "  Terminal Layout Demo  "
	rootStyled.AddChild(headerStyled)

	// Simple message
	msg := components.NewMessageBlock("Hello from the layout engine!\nThis is a test of the terminal UI system.")
	rootStyled.AddChild(msg.ToStyledNode())

	// Footer
	footerNode := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   float64(width - 4),
			Height:  3,
		},
	}
	fgGray, _ := color.ParseColor("#888888")
	footerStyle := &renderer.Style{
		Foreground: &fgGray,
		Dim:        true,
	}
	footerStyled := renderer.NewStyledNode(footerNode, footerStyle)
	footerStyled.Content = fmt.Sprintf("Terminal: %dx%d", width, height)
	rootStyled.AddChild(footerStyled)

	// Perform layout calculation
	constraints := layout.Tight(float64(width), float64(height))
	size := layout.Layout(root, constraints)

	fmt.Printf("Layout calculated: %v x %v\n", size.Width, size.Height)
	fmt.Printf("Root rect: X=%v Y=%v W=%v H=%v\n",
		root.Rect.X, root.Rect.Y, root.Rect.Width, root.Rect.Height)

	if len(root.Children) > 0 {
		fmt.Printf("First child rect: X=%v Y=%v W=%v H=%v\n",
			root.Children[0].Rect.X, root.Children[0].Rect.Y,
			root.Children[0].Rect.Width, root.Children[0].Rect.Height)
	}

	fmt.Println("\n--- Rendered Output ---\n")

	// Render to screen
	screen.Render(rootStyled)

	// Output
	fmt.Print(screen.String())
}
