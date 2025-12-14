package main

import (
	"fmt"

	"github.com/SCKelemen/color"
	"github.com/SCKelemen/layout"
	"github.com/SCKelemen/cli/renderer"
)

func main() {
	// Detect terminal capabilities
	caps := renderer.DetectCapabilities()

	fmt.Println("=== Terminal Capability Detection ===")
	fmt.Printf("Color Mode: %s\n", caps.ColorMode)
	fmt.Printf("Is TTY: %t\n", caps.IsTTY)
	fmt.Println()

	if !caps.IsTTY {
		fmt.Println("⚠️  Not running in a terminal. Colors may not display correctly.")
		fmt.Println()
	}

	// Show gradients in all supported modes
	fmt.Println("=== Color Degradation Demo ===")
	fmt.Println("Showing the same gradient in different color modes:\n")

	// True Color
	fmt.Println("1. True Color (24-bit) - What the gradient should look like:")
	showGradient(renderer.ColorModeTrueColor, 70)

	// 256 Color
	fmt.Println("\n2. 256 Color Mode - Approximation using 256-color palette:")
	showGradient(renderer.ColorMode256, 70)

	// 16 Color
	fmt.Println("\n3. 16 Color Mode - Basic ANSI colors:")
	showGradient(renderer.ColorMode16, 70)

	// Current terminal
	fmt.Printf("\n4. Your Terminal (%s):\n", caps.ColorMode)
	showGradient(caps.ColorMode, 70)

	// Color swatches
	fmt.Println("\n=== Color Palette ===")
	showColorSwatches(caps.ColorMode)

	fmt.Println("\n=== Progressive Enhancement ===")
	fmt.Println("✓ Works in all terminals (graceful degradation)")
	fmt.Println("✓ Better experience in modern terminals")
	fmt.Println("✓ Falls back to 256 colors, 16 colors, or plain text")
	fmt.Println()

	if caps.ColorMode == renderer.ColorMode256 {
		fmt.Println("💡 Your terminal supports 256 colors!")
		fmt.Println("   For even better colors, try: iTerm2, Alacritty, or Kitty")
		fmt.Println("   Or set: export COLORTERM=truecolor")
	}
}

func showGradient(mode renderer.ColorMode, width int) {
	screen := renderer.NewScreen(width, 2)

	root := &layout.Node{
		Style: layout.Style{
			Display:       layout.DisplayFlex,
			FlexDirection: layout.FlexDirectionRow,
			Width:         float64(width),
			Height:        2,
		},
	}
	rootStyled := renderer.NewStyledNode(root, nil)

	// Create gradient cells (Blue to Red)
	for i := 0; i < width; i++ {
		t := float64(i) / float64(width-1)
		hue := 240 + (t * 150) // Blue to Red
		if hue > 360 {
			hue -= 360
		}
		colorStr := fmt.Sprintf("oklch(0.6 0.2 %.0f)", hue)
		c, _ := color.ParseColor(colorStr)

		cellNode := &layout.Node{
			Style: layout.Style{
				Display: layout.DisplayBlock,
				Width:   1,
				Height:  2,
			},
		}
		cellStyle := &renderer.Style{
			Background: &c,
		}
		cellStyled := renderer.NewStyledNode(cellNode, cellStyle)
		cellStyled.Content = " "
		rootStyled.AddChild(cellStyled)
	}

	constraints := layout.Tight(float64(width), 2)
	layout.Layout(root, constraints)

	// Override the renderer's color mode
	screen.SetColorMode(mode)
	screen.Render(rootStyled)
	fmt.Print(screen.String())
}

func showColorSwatches(mode renderer.ColorMode) {
	swatches := []struct {
		name  string
		oklch string
	}{
		{"Red", "oklch(0.55 0.22 30)"},
		{"Orange", "oklch(0.65 0.20 60)"},
		{"Yellow", "oklch(0.80 0.15 90)"},
		{"Green", "oklch(0.60 0.18 140)"},
		{"Cyan", "oklch(0.70 0.15 200)"},
		{"Blue", "oklch(0.50 0.20 260)"},
		{"Purple", "oklch(0.55 0.20 300)"},
		{"Magenta", "oklch(0.60 0.25 330)"},
	}

	width := 70
	screen := renderer.NewScreen(width, 4)

	root := &layout.Node{
		Style: layout.Style{
			Display:       layout.DisplayFlex,
			FlexDirection: layout.FlexDirectionRow,
			Width:         float64(width),
			Height:        4,
		},
	}
	rootStyled := renderer.NewStyledNode(root, nil)

	swatchWidth := 8

	for _, swatch := range swatches {
		c, _ := color.ParseColor(swatch.oklch)
		fgWhite, _ := color.ParseColor("#FFFFFF")

		swatchNode := &layout.Node{
			Style: layout.Style{
				Display: layout.DisplayBlock,
				Width:   float64(swatchWidth),
				Height:  4,
				Margin:  layout.Spacing{Top: 0, Right: 1, Bottom: 0, Left: 0},
			},
		}
		swatchStyle := &renderer.Style{
			Background:  &c,
			Foreground:  &fgWhite,
			Bold:        true,
			BorderColor: &c,
		}
		swatchStyle.WithBorder(renderer.RoundedBorder)
		swatchStyled := renderer.NewStyledNode(swatchNode, swatchStyle)
		swatchStyled.Content = fmt.Sprintf("\n %s", swatch.name)
		rootStyled.AddChild(swatchStyled)
	}

	constraints := layout.Tight(float64(width), 4)
	layout.Layout(root, constraints)

	screen.SetColorMode(mode)
	screen.Render(rootStyled)
	fmt.Print(screen.String())
}
