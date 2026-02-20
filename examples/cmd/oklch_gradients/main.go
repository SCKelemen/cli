package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/SCKelemen/cli/renderer"
	"github.com/SCKelemen/color"
	"github.com/SCKelemen/layout"
)

func main() {
	// Clear screen
	fmt.Print("\033[2J\033[H")

	fmt.Println("╔═══════════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║               OKLCH Color Gradients - Terminal Showcase                  ║")
	fmt.Print("╚═══════════════════════════════════════════════════════════════════════════╝\n")

	// Rainbow Hue Gradient
	fmt.Println("1. Hue Rotation (0° → 360°) - Full Color Wheel:")
	printGradient(70, func(t float64) string {
		hue := t * 360
		return fmt.Sprintf("oklch(0.65 0.2 %.0f)", hue)
	})
	fmt.Println()

	// Lightness Gradient
	fmt.Println("2. Lightness Variation (Dark → Light):")
	printGradient(70, func(t float64) string {
		lightness := 0.2 + (t * 0.7)
		return fmt.Sprintf("oklch(%.2f 0.15 270)", lightness)
	})
	fmt.Println()

	// Chroma Gradient
	fmt.Println("3. Chroma/Saturation (Gray → Vivid):")
	printGradient(70, func(t float64) string {
		chroma := t * 0.3
		return fmt.Sprintf("oklch(0.6 %.2f 180)", chroma)
	})
	fmt.Println()

	// Blue to Red
	fmt.Println("4. Blue → Purple → Red (Smooth OKLCH Transition):")
	printGradient(70, func(t float64) string {
		hue := 240 + (t * 150)
		if hue > 360 {
			hue -= 360
		}
		return fmt.Sprintf("oklch(0.6 0.2 %.0f)", hue)
	})
	fmt.Println()

	// Sunset Gradient
	fmt.Println("5. Sunset Gradient (Purple → Pink → Orange → Yellow):")
	printGradient(70, func(t float64) string {
		if t < 0.33 {
			hue := 280 + (t/0.33)*40
			return fmt.Sprintf("oklch(0.5 0.25 %.0f)", hue)
		} else if t < 0.66 {
			hue := 320 + ((t-0.33)/0.33)*30
			return fmt.Sprintf("oklch(0.65 0.2 %.0f)", hue)
		} else {
			hue := 50 + ((t-0.66)/0.34)*30
			return fmt.Sprintf("oklch(0.75 0.18 %.0f)", hue)
		}
	})
	fmt.Println()

	// Ocean Gradient
	fmt.Println("6. Ocean Gradient (Deep Blue → Cyan → Teal):")
	printGradient(70, func(t float64) string {
		hue := 220 + (t * 50)
		lightness := 0.4 + (t * 0.3)
		return fmt.Sprintf("oklch(%.2f 0.15 %.0f)", lightness, hue)
	})
	fmt.Println()

	// Warm to Cool
	fmt.Println("7. Warm → Cool (Red → Orange → Green → Blue):")
	printGradient(70, func(t float64) string {
		hue := 30 + (t * 210)
		return fmt.Sprintf("oklch(0.6 0.18 %.0f)", hue)
	})
	fmt.Println()

	// Color Swatches
	fmt.Print("\n8. Color Palette (Perceptually Uniform):")
	printColorSwatches()

	fmt.Print("\n" + renderColorInfo())

	// Wait for user
	fmt.Print("\nPress Enter to exit...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func printGradient(width int, colorFunc func(float64) string) {
	screen := renderer.NewScreen(width, 2)

	root := &layout.Node{
		Style: layout.Style{
			Display:       layout.DisplayFlex,
			FlexDirection: layout.FlexDirectionRow,
			Width:         layout.Px(float64(width)),
			Height:        layout.Px(2),
		},
	}
	rootStyled := renderer.NewStyledNode(root, nil)

	// Create gradient cells
	for i := 0; i < width; i++ {
		t := float64(i) / float64(width-1)
		colorStr := colorFunc(t)
		c, _ := color.ParseColor(colorStr)

		cellNode := &layout.Node{
			Style: layout.Style{
				Display: layout.DisplayBlock,
				Width:   layout.Px(1),
				Height:  layout.Px(2),
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
	ctx := &layout.LayoutContext{
		ViewportWidth:  float64(width),
		ViewportHeight: 2,
		RootFontSize:   16,
	}
	layout.Layout(root, constraints, ctx)
	screen.Render(rootStyled)
	fmt.Print(screen.String())
}

func printColorSwatches() {
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

	// 8 swatches * 10 width = 80 columns total
	width := 80
	screen := renderer.NewScreen(width, 4)

	root := &layout.Node{
		Style: layout.Style{
			Display:       layout.DisplayFlex,
			FlexDirection: layout.FlexDirectionRow,
			Width:         layout.Px(float64(width)),
			Height:        layout.Px(4),
		},
	}
	rootStyled := renderer.NewStyledNode(root, nil)

	swatchWidth := 10 // Width 10 allows 8 content columns (fits " Magenta")

	for _, swatch := range swatches {
		c, _ := color.ParseColor(swatch.oklch)
		fgWhite, _ := color.ParseColor("#FFFFFF")

		swatchNode := &layout.Node{
			Style: layout.Style{
				Display: layout.DisplayBlock,
				Width:   layout.Px(float64(swatchWidth)),
				Height:  layout.Px(4),
				Margin:  layout.Spacing{Top: layout.Px(0), Right: layout.Px(0), Bottom: layout.Px(0), Left: layout.Px(0)}, // No margins, boxes are adjacent
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
	ctx := &layout.LayoutContext{
		ViewportWidth:  float64(width),
		ViewportHeight: 4,
		RootFontSize:   16,
	}
	layout.Layout(root, constraints, ctx)
	screen.Render(rootStyled)
	fmt.Print(screen.String())
}

func renderColorInfo() string {
	return `
┌────────────────────────────────────────────────────────────────────────┐
│ About OKLCH Colors                                                     │
├────────────────────────────────────────────────────────────────────────┤
│ • L (Lightness): 0.0 (black) to 1.0 (white)                         │
│ • C (Chroma): 0.0 (gray) to ~0.4 (vivid)                            │
│ • H (Hue): 0° to 360° (color wheel)                                 │
│                                                                        │
│ Benefits:                                                              │
│ ✓ Perceptually uniform - equal changes look equal to the human eye  │
│ ✓ Smooth gradients - no muddy middle colors like RGB                │
│ ✓ Predictable lightness - same L value = same perceived brightness  │
│ ✓ Accessible color systems - easier to maintain contrast ratios     │
└────────────────────────────────────────────────────────────────────────┘
`
}
