package main

import (
	"fmt"

	"github.com/SCKelemen/cli/renderer"
	"github.com/SCKelemen/color"
)

func main() {
	fmt.Print("=== Color Conversion Debug ===\n")

	// Test OKLCH color conversion
	oklchColors := []string{
		"oklch(0.55 0.22 30)",  // Red
		"oklch(0.65 0.20 60)",  // Orange
		"oklch(0.80 0.15 90)",  // Yellow
		"oklch(0.60 0.18 140)", // Green
		"oklch(0.70 0.15 200)", // Cyan
		"oklch(0.50 0.20 260)", // Blue
		"oklch(0.55 0.20 300)", // Purple
		"oklch(0.60 0.25 330)", // Magenta
	}

	caps := renderer.DetectCapabilities()
	fmt.Printf("Detected mode: %s\n\n", caps.ColorMode)

	for _, oklch := range oklchColors {
		c, err := color.ParseColor(oklch)
		if err != nil {
			fmt.Printf("Error parsing %s: %v\n", oklch, err)
			continue
		}

		r, g, b, _ := c.RGBA()
		r8 := int(r * 255)
		g8 := int(g * 255)
		b8 := int(b * 255)

		fmt.Printf("%s → RGB(%d,%d,%d)\n", oklch, r8, g8, b8)

		// Test with 256-color mode
		ansiRenderer := renderer.NewANSIRendererWithMode(renderer.ColorMode256)
		style := &renderer.Style{Background: &c}
		fmt.Print("  256-color: ")
		fmt.Print(ansiRenderer.RenderStyle(style))
		fmt.Print("  SAMPLE  ")
		fmt.Print(ansiRenderer.Reset())
		fmt.Println()

		// Test with 16-color mode
		ansiRenderer16 := renderer.NewANSIRendererWithMode(renderer.ColorMode16)
		style16 := &renderer.Style{Background: &c}
		fmt.Print("  16-color:  ")
		fmt.Print(ansiRenderer16.RenderStyle(style16))
		fmt.Print("  SAMPLE  ")
		fmt.Print(ansiRenderer16.Reset())
		fmt.Print("\n")
	}
}
