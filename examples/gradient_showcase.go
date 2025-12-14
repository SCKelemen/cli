package main

import (
	"fmt"
	"strings"

	"github.com/SCKelemen/color"
	"github.com/SCKelemen/layout"
	"github.com/SCKelemen/cli/renderer"
)

func main() {
	fmt.Println("=== OKLCH Gradient Showcase ===\n")

	width, height := 90, 50
	screen := renderer.NewScreen(width, height)
	root := buildGradientShowcase(width, height)

	constraints := layout.Tight(float64(width), float64(height))
	layout.Layout(root.Node, constraints)
	screen.Render(root)

	fmt.Print(screen.String())
}

func buildGradientShowcase(width, height int) *renderer.StyledNode {
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
	header := createHeader(width-4, "OKLCH Color Gradients")
	rootStyled.AddChild(header)

	// Gradient 1: Hue Spectrum (Rainbow)
	addLabel(rootStyled, "1. Hue Rotation (0° → 360°) - Full Color Wheel", width-4)
	gradient1 := createGradient(width-4, 3, func(t float64) color.Color {
		hue := t * 360
		c, _ := color.ParseColor(fmt.Sprintf("oklch(0.65 0.2 %.0f)", hue))
		return c
	})
	rootStyled.AddChild(gradient1)

	// Gradient 2: Lightness
	addLabel(rootStyled, "2. Lightness Variation (Dark → Light)", width-4)
	gradient2 := createGradient(width-4, 3, func(t float64) color.Color {
		lightness := 0.2 + (t * 0.7)
		c, _ := color.ParseColor(fmt.Sprintf("oklch(%.2f 0.15 270)", lightness))
		return c
	})
	rootStyled.AddChild(gradient2)

	// Gradient 3: Chroma (Saturation)
	addLabel(rootStyled, "3. Chroma/Saturation (Gray → Vivid)", width-4)
	gradient3 := createGradient(width-4, 3, func(t float64) color.Color {
		chroma := t * 0.3
		c, _ := color.ParseColor(fmt.Sprintf("oklch(0.6 %.2f 180)", chroma))
		return c
	})
	rootStyled.AddChild(gradient3)

	// Gradient 4: Blue to Red (smooth transition)
	addLabel(rootStyled, "4. Blue → Purple → Red (Smooth OKLCH interpolation)", width-4)
	gradient4 := createGradient(width-4, 3, func(t float64) color.Color {
		hue := 240 + (t * 150) // 240° (blue) to 390° (red wrapping around)
		if hue > 360 {
			hue -= 360
		}
		c, _ := color.ParseColor(fmt.Sprintf("oklch(0.6 0.2 %.0f)", hue))
		return c
	})
	rootStyled.AddChild(gradient4)

	// Gradient 5: Sunset
	addLabel(rootStyled, "5. Sunset Gradient (Purple → Pink → Orange → Yellow)", width-4)
	gradient5 := createGradient(width-4, 3, func(t float64) color.Color {
		if t < 0.33 {
			hue := 280 + (t/0.33)*40
			c, _ := color.ParseColor(fmt.Sprintf("oklch(0.5 0.25 %.0f)", hue))
			return c
		} else if t < 0.66 {
			hue := 320 + ((t-0.33)/0.33)*30
			c, _ := color.ParseColor(fmt.Sprintf("oklch(0.65 0.2 %.0f)", hue))
			return c
		} else {
			hue := 50 + ((t-0.66)/0.34)*30
			c, _ := color.ParseColor(fmt.Sprintf("oklch(0.75 0.18 %.0f)", hue))
			return c
		}
	})
	rootStyled.AddChild(gradient5)

	// Gradient 6: Ocean
	addLabel(rootStyled, "6. Ocean Gradient (Deep Blue → Cyan → Teal)", width-4)
	gradient6 := createGradient(width-4, 3, func(t float64) color.Color {
		hue := 220 + (t * 50)
		lightness := 0.4 + (t * 0.3)
		c, _ := color.ParseColor(fmt.Sprintf("oklch(%.2f 0.15 %.0f)", lightness, hue))
		return c
	})
	rootStyled.AddChild(gradient6)

	// Footer
	footer := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   float64(width - 4),
			Height:  3,
			Margin:  layout.Spacing{Top: 1, Right: 0, Bottom: 0, Left: 0},
		},
	}
	fgGray, _ := color.ParseColor("#888888")
	footerStyle := &renderer.Style{
		Foreground: &fgGray,
		Dim:        true,
	}
	footerStyled := renderer.NewStyledNode(footer, footerStyle)
	footerStyled.Content = "OKLCH (Lightness, Chroma, Hue) provides perceptually uniform gradients\nColors appear to transition smoothly to the human eye, unlike RGB"
	rootStyled.AddChild(footerStyled)

	return rootStyled
}

func createHeader(width int, title string) *renderer.StyledNode {
	node := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   float64(width),
			Height:  5,
		},
	}
	fgWhite, _ := color.ParseColor("#FFFFFF")
	bgPurple, _ := color.ParseColor("oklch(0.5 0.2 270)")
	borderPurple, _ := color.ParseColor("oklch(0.7 0.2 270)")
	style := &renderer.Style{
		Foreground:  &fgWhite,
		Background:  &bgPurple,
		Bold:        true,
		BorderColor: &borderPurple,
	}
	style.WithBorder(renderer.DoubleBorder)
	styled := renderer.NewStyledNode(node, style)
	styled.Content = fmt.Sprintf("  %s  ", title)
	return styled
}

func addLabel(parent *renderer.StyledNode, text string, width int) {
	node := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   float64(width),
			Height:  1,
			Margin:  layout.Spacing{Top: 1, Right: 0, Bottom: 0, Left: 0},
		},
	}
	fgLabel, _ := color.ParseColor("#CCCCCC")
	style := &renderer.Style{
		Foreground: &fgLabel,
	}
	styled := renderer.NewStyledNode(node, style)
	styled.Content = text
	parent.AddChild(styled)
}

func createGradient(width int, height int, colorFunc func(float64) color.Color) *renderer.StyledNode {
	node := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   float64(width),
			Height:  float64(height),
		},
	}

	// Generate gradient as a string of colored blocks
	var content strings.Builder
	steps := width - 2 // Account for borders

	for y := 0; y < height; y++ {
		for x := 0; x < steps; x++ {
			t := float64(x) / float64(steps-1)
			_ = colorFunc(t)
			content.WriteRune('█')
		}
		if y < height-1 {
			content.WriteRune('\n')
		}
	}

	// For now, use a single color for the gradient bar
	// The ideal implementation would support per-character coloring
	midColor := colorFunc(0.5)
	style := &renderer.Style{
		Background: &midColor,
	}

	styled := renderer.NewStyledNode(node, style)
	styled.Content = content.String()
	return styled
}
