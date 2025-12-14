package main

import (
	"fmt"

	"github.com/SCKelemen/color"
	"github.com/SCKelemen/layout"
	"github.com/SCKelemen/cli/renderer"
)

func main() {
	fmt.Println("=== OKLCH Color Gradients Demo ===\n")

	// Terminal dimensions
	width, height := 100, 45

	// Create screen
	screen := renderer.NewScreen(width, height)

	// Build color demo layout
	root := buildColorDemo(width, height)

	// Perform layout calculation
	constraints := layout.Tight(float64(width), float64(height))
	layout.Layout(root.Node, constraints)

	// Render to screen
	screen.Render(root)

	// Output
	fmt.Print(screen.String())
}

func buildColorDemo(width, height int) *renderer.StyledNode {
	// Create root container
	root := &layout.Node{
		Style: layout.Style{
			Display:       layout.DisplayFlex,
			FlexDirection: layout.FlexDirectionColumn,
			Width:         float64(width),
			Height:        float64(height),
			Padding:       layout.Spacing{Top: 1, Right: 2, Bottom: 1, Left: 2},
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
	headerStyled.Content = "  OKLCH Gradient Showcase  "
	rootStyled.AddChild(headerStyled)

	// Gradient 1: Hue rotation (rainbow)
	addGradientBar(rootStyled, "Hue Rotation (0° → 360°)", width-4, func(t float64) color.Color {
		hue := t * 360
		c, _ := color.ParseColor(fmt.Sprintf("oklch(0.65 0.2 %.0f)", hue))
		return c
	})

	// Gradient 2: Lightness variation
	addGradientBar(rootStyled, "Lightness (0.2 → 0.9)", width-4, func(t float64) color.Color {
		lightness := 0.2 + (t * 0.7)
		c, _ := color.ParseColor(fmt.Sprintf("oklch(%.2f 0.15 270)", lightness))
		return c
	})

	// Gradient 3: Chroma variation (saturation)
	addGradientBar(rootStyled, "Chroma/Saturation (0.0 → 0.3)", width-4, func(t float64) color.Color {
		chroma := t * 0.3
		c, _ := color.ParseColor(fmt.Sprintf("oklch(0.6 %.2f 180)", chroma))
		return c
	})

	// Gradient 4: Blue to Red (smooth OKLCH interpolation)
	addGradientBar(rootStyled, "Blue → Red (OKLCH)", width-4, func(t float64) color.Color {
		startHue := 240.0
		endHue := 30.0
		hue := startHue + (t * (endHue - startHue))
		c, _ := color.ParseColor(fmt.Sprintf("oklch(0.6 0.2 %.0f)", hue))
		return c
	})

	// Gradient 5: Multi-stop gradient
	addGradientBar(rootStyled, "Sunset Gradient", width-4, func(t float64) color.Color {
		if t < 0.33 {
			// Deep purple to pink
			hue := 280 + (t / 0.33 * 40)
			c, _ := color.ParseColor(fmt.Sprintf("oklch(0.5 0.25 %.0f)", hue))
			return c
		} else if t < 0.66 {
			// Pink to orange
			hue := 320 + ((t - 0.33) / 0.33 * 30)
			c, _ := color.ParseColor(fmt.Sprintf("oklch(0.65 0.2 %.0f)", hue))
			return c
		} else {
			// Orange to yellow
			hue := 50 + ((t - 0.66) / 0.34 * 30)
			c, _ := color.ParseColor(fmt.Sprintf("oklch(0.75 0.18 %.0f)", hue))
			return c
		}
	})

	// Color swatches
	addColorSwatches(rootStyled, width-4)

	// Footer
	footerNode := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   float64(width - 4),
			Height:  2,
		},
	}
	fgGray, _ := color.ParseColor("#888888")
	footerStyle := &renderer.Style{
		Foreground: &fgGray,
		Dim:        true,
	}
	footerStyled := renderer.NewStyledNode(footerNode, footerStyle)
	footerStyled.Content = "OKLCH provides perceptually uniform color interpolation\nPerfect for smooth gradients and accessible color systems"
	rootStyled.AddChild(footerStyled)

	return rootStyled
}

func addGradientBar(parent *renderer.StyledNode, label string, width int, colorFunc func(float64) color.Color) {
	// Label
	labelNode := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   float64(width),
			Height:  1,
			Margin:  layout.Spacing{Top: 1, Right: 0, Bottom: 0, Left: 0},
		},
	}
	fgLabel, _ := color.ParseColor("#AAAAAA")
	labelStyle := &renderer.Style{
		Foreground: &fgLabel,
	}
	labelStyled := renderer.NewStyledNode(labelNode, labelStyle)
	labelStyled.Content = label
	parent.AddChild(labelStyled)

	// Gradient bar container
	barContainer := &layout.Node{
		Style: layout.Style{
			Display:       layout.DisplayFlex,
			FlexDirection: layout.FlexDirectionRow,
			Width:         float64(width),
			Height:        3,
		},
	}
	barContainerStyled := renderer.NewStyledNode(barContainer, nil)

	// Create gradient cells
	numCells := width - 2
	for i := 0; i < numCells; i++ {
		t := float64(i) / float64(numCells-1)
		c := colorFunc(t)

		cellNode := &layout.Node{
			Style: layout.Style{
				Display: layout.DisplayBlock,
				Width:   1,
				Height:  3,
			},
		}
		cellStyle := &renderer.Style{
			Background: &c,
		}
		cellStyled := renderer.NewStyledNode(cellNode, cellStyle)
		cellStyled.Content = " "
		barContainerStyled.AddChild(cellStyled)
	}

	parent.AddChild(barContainerStyled)
}

func addColorSwatches(parent *renderer.StyledNode, width int) {
	// Section label
	labelNode := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   float64(width),
			Height:  1,
			Margin:  layout.Spacing{Top: 1, Right: 0, Bottom: 0, Left: 0},
		},
	}
	fgLabel, _ := color.ParseColor("#AAAAAA")
	labelStyle := &renderer.Style{
		Foreground: &fgLabel,
	}
	labelStyled := renderer.NewStyledNode(labelNode, labelStyle)
	labelStyled.Content = "Color Palette Examples"
	parent.AddChild(labelStyled)

	// Swatches container
	swatchContainer := &layout.Node{
		Style: layout.Style{
			Display:       layout.DisplayFlex,
			FlexDirection: layout.FlexDirectionRow,
			Width:         float64(width),
			Height:        5,
		},
	}
	swatchContainerStyled := renderer.NewStyledNode(swatchContainer, nil)

	// Define color swatches
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
		{"Pink", "oklch(0.65 0.20 340)"},
	}

	swatchWidth := (width - len(swatches) - 2) / len(swatches)

	for _, swatch := range swatches {
		c, _ := color.ParseColor(swatch.oklch)
		fgWhite, _ := color.ParseColor("#FFFFFF")

		swatchNode := &layout.Node{
			Style: layout.Style{
				Display: layout.DisplayBlock,
				Width:   float64(swatchWidth),
				Height:  5,
				Margin:  layout.Spacing{Top: 0, Right: 1, Bottom: 0, Left: 0},
			},
		}
		swatchStyle := &renderer.Style{
			Background:  &c,
			Foreground:  &fgWhite,
			BorderColor: &c,
			Bold:        true,
		}
		swatchStyle.WithBorder(renderer.RoundedBorder)
		swatchStyled := renderer.NewStyledNode(swatchNode, swatchStyle)
		swatchStyled.Content = fmt.Sprintf("\n %s", swatch.name)
		swatchContainerStyled.AddChild(swatchStyled)
	}

	parent.AddChild(swatchContainerStyled)
}
