package main

import (
	"fmt"

	"github.com/SCKelemen/color"
)

func main() {
	fmt.Print("=== OKLCH to RGB Conversion Test ===\n")

	// Test a simple gradient like in oklch_gradients.go
	fmt.Println("Blue to Red gradient (hue 240 to 30):")
	for i := 0; i < 10; i++ {
		t := float64(i) / 9.0
		hue := 240 + (t * 150)
		if hue > 360 {
			hue -= 360
		}
		colorStr := fmt.Sprintf("oklch(0.6 0.2 %.0f)", hue)
		c, err := color.ParseColor(colorStr)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		r, g, b, _ := c.RGBA()
		r8 := int(r * 255)
		g8 := int(g * 255)
		b8 := int(b * 255)

		fmt.Printf("%d. %s → RGB(%3d,%3d,%3d) ", i, colorStr, r8, g8, b8)

		// Show as 256-color background
		fmt.Printf("\x1b[48;2;%d;%d;%dm     \x1b[0m ", r8, g8, b8)

		// Check if it's grayscale (all RGB values similar)
		maxDiff := max(abs(r8-g8), max(abs(g8-b8), abs(r8-b8)))
		if maxDiff < 20 {
			fmt.Print("⚠️  GRAYSCALE!")
		}
		fmt.Println()
	}

	fmt.Print("\n\nRainbow gradient (hue 0 to 360):")
	for i := 0; i < 10; i++ {
		t := float64(i) / 9.0
		hue := t * 360
		colorStr := fmt.Sprintf("oklch(0.65 0.2 %.0f)", hue)
		c, err := color.ParseColor(colorStr)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		r, g, b, _ := c.RGBA()
		r8 := int(r * 255)
		g8 := int(g * 255)
		b8 := int(b * 255)

		fmt.Printf("%d. %s → RGB(%3d,%3d,%3d) ", i, colorStr, r8, g8, b8)
		fmt.Printf("\x1b[48;2;%d;%d;%dm     \x1b[0m ", r8, g8, b8)

		maxDiff := max(abs(r8-g8), max(abs(g8-b8), abs(r8-b8)))
		if maxDiff < 20 {
			fmt.Print("⚠️  GRAYSCALE!")
		}
		fmt.Println()
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
