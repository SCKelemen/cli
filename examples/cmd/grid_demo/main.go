package main

import (
	"fmt"

	"github.com/SCKelemen/cli/renderer"
	"github.com/SCKelemen/color"
	"github.com/SCKelemen/layout"
)

func main() {
	fmt.Print("=== Grid Layout Demo ===\n")

	// Terminal dimensions
	width, height := 100, 40

	// Create screen
	screen := renderer.NewScreen(width, height)

	// Build grid layout
	root := buildGridLayout(width, height)

	// Perform layout calculation
	constraints := layout.Tight(float64(width), float64(height))
	ctx := &layout.LayoutContext{
		ViewportWidth:  float64(width),
		ViewportHeight: float64(height),
		RootFontSize:   16,
	}
	size := layout.Layout(root.Node, constraints, ctx)

	fmt.Printf("Layout: %.0fx%.0f\n\n", size.Width, size.Height)

	// Render to screen
	screen.Render(root)

	// Output
	fmt.Print(screen.String())
}

func buildGridLayout(width, height int) *renderer.StyledNode {
	// Create root container
	root := &layout.Node{
		Style: layout.Style{
			Display:       layout.DisplayFlex,
			FlexDirection: layout.FlexDirectionColumn,
			Width:         layout.Px(float64(width)),
			Height:        layout.Px(float64(height)),
			Padding:       layout.Spacing{Top: layout.Px(2), Right: layout.Px(2), Bottom: layout.Px(2), Left: layout.Px(2)},
		},
	}
	rootStyled := renderer.NewStyledNode(root, nil)

	// Header
	headerNode := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   layout.Px(float64(width - 4)),
			Height:  layout.Px(5),
		},
	}
	fgWhite, _ := color.ParseColor("#FFFFFF")
	bgBlue, _ := color.ParseColor("oklch(0.5 0.2 240)")
	borderBlue, _ := color.ParseColor("oklch(0.7 0.2 240)")
	headerStyle := &renderer.Style{
		Foreground:  &fgWhite,
		Background:  &bgBlue,
		Bold:        true,
		BorderColor: &borderBlue,
	}
	headerStyle.WithBorder(renderer.ThickBorder)
	headerStyled := renderer.NewStyledNode(headerNode, headerStyle)
	headerStyled.Content = "  CSS Grid Layout - Dashboard Example  "
	rootStyled.AddChild(headerStyled)

	// Grid container
	gridNode := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayGrid,
			// 3 columns: 30%, 40%, 30%
			GridTemplateColumns: []layout.GridTrack{
				layout.FractionTrack(3),
				layout.FractionTrack(4),
				layout.FractionTrack(3),
			},
			// 3 rows: auto-sized
			GridTemplateRows: []layout.GridTrack{
				layout.FractionTrack(1),
				layout.FractionTrack(1),
				layout.FractionTrack(1),
			},
			GridGap: layout.Px(2),
			Width:   layout.Px(float64(width - 4)),
			Height:  layout.Px(float64(height - 10)),
		},
	}
	gridStyled := renderer.NewStyledNode(gridNode, nil)

	// Card colors
	colors := []struct {
		bg     string
		border string
		title  string
	}{
		{"oklch(0.45 0.15 30)", "oklch(0.65 0.15 30)", "System Stats"},
		{"oklch(0.45 0.15 90)", "oklch(0.65 0.15 90)", "Network"},
		{"oklch(0.45 0.15 150)", "oklch(0.65 0.15 150)", "CPU Usage"},
		{"oklch(0.45 0.15 210)", "oklch(0.65 0.15 210)", "Memory"},
		{"oklch(0.45 0.15 270)", "oklch(0.65 0.15 270)", "Processes"},
		{"oklch(0.45 0.15 330)", "oklch(0.65 0.15 330)", "Disk I/O"},
		{"oklch(0.45 0.15 0)", "oklch(0.65 0.15 0)", "Logs"},
		{"oklch(0.45 0.15 60)", "oklch(0.65 0.15 60)", "Alerts"},
		{"oklch(0.45 0.15 120)", "oklch(0.65 0.15 120)", "Status"},
	}

	// Create 9 cards in the grid
	for i := 0; i < 9; i++ {
		cardNode := &layout.Node{
			Style: layout.Style{
				Display: layout.DisplayBlock,
			},
		}

		colorSet := colors[i%len(colors)]
		bg, _ := color.ParseColor(colorSet.bg)
		border, _ := color.ParseColor(colorSet.border)
		fgCard, _ := color.ParseColor("#FFFFFF")

		cardStyle := &renderer.Style{
			Foreground:  &fgCard,
			Background:  &bg,
			BorderColor: &border,
			Bold:        true,
		}
		cardStyle.WithBorder(renderer.RoundedBorder)

		cardStyled := renderer.NewStyledNode(cardNode, cardStyle)
		cardStyled.Content = fmt.Sprintf(" %s\n\n Value: %d\n Status: OK", colorSet.title, (i+1)*123)

		gridStyled.AddChild(cardStyled)
	}

	rootStyled.AddChild(gridStyled)

	// Footer with grid info
	footerNode := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   layout.Px(float64(width - 4)),
			Height:  layout.Px(2),
		},
	}
	fgGray, _ := color.ParseColor("#888888")
	footerStyle := &renderer.Style{
		Foreground: &fgGray,
		Dim:        true,
	}
	footerStyled := renderer.NewStyledNode(footerNode, footerStyle)
	footerStyled.Content = "Grid: 3 columns (3fr, 4fr, 3fr) × 3 rows (1fr each) with 2-unit gap"
	rootStyled.AddChild(footerStyled)

	return rootStyled
}
