package main

import (
	"fmt"

	"github.com/SCKelemen/cli/renderer"
	"github.com/SCKelemen/color"
	"github.com/SCKelemen/layout"
)

func main() {
	fmt.Print("\033[2J\033[H") // Clear screen

	width := 76
	height := 8

	screen := renderer.NewScreen(width, height)

	root := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   layout.Px(float64(width)),
			Height:  layout.Px(float64(height)),
		},
	}

	borderColor, _ := color.ParseColor("#7D56F4")
	textColor, _ := color.ParseColor("#FAFAFA")

	style := &renderer.Style{
		Foreground:  &textColor,
		BorderColor: &borderColor,
	}
	style.WithBorder(renderer.RoundedBorder)

	rootStyled := renderer.NewStyledNode(root, style)
	rootStyled.Content = `
Test without special chars
Another regular line here
• Line with bullet point
✓ Line with checkmark
Final regular line`

	constraints := layout.Tight(float64(width), float64(height))
	ctx := &layout.LayoutContext{
		ViewportWidth:  float64(width),
		ViewportHeight: float64(height),
		RootFontSize:   16,
	}
	layout.Layout(root, constraints, ctx)
	screen.Render(rootStyled)

	output := screen.String()
	fmt.Print(output)

	// Debug: show what we rendered character by character
	fmt.Print("\n\nDebug - character positions:")
	lines := screen.Cells
	for i, line := range lines {
		fmt.Printf("Line %d (len=%d): ", i, len(line))
		for j, cell := range line {
			if cell.Content == "│" || cell.Content == "┐" || cell.Content == "┘" {
				fmt.Printf("[border@%d]", j)
			}
		}
		fmt.Println()
	}

	fmt.Print("\nPress Enter to exit...")
	fmt.Scanln()
}
