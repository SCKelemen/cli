package main

import (
	"fmt"

	"github.com/SCKelemen/color"
	"github.com/SCKelemen/layout"
	"github.com/SCKelemen/cli/renderer"
)

func main() {
	width, height := 100, 30

	screen := renderer.NewScreen(width, height)
	root := buildSimpleDashboard(width, height)

	constraints := layout.Tight(float64(width), float64(height))
	ctx := &layout.LayoutContext{
		ViewportWidth:  float64(width),
		ViewportHeight: float64(height),
		RootFontSize:   16,
	}
	layout.Layout(root.Node, constraints, ctx)
	screen.Render(root)

	fmt.Print(screen.String())
}

func buildSimpleDashboard(width, height int) *renderer.StyledNode {
	root := &layout.Node{
		Style: layout.Style{
			Display:       layout.DisplayFlex,
			FlexDirection: layout.FlexDirectionColumn,
			Width:         layout.Px(float64(width)),
			Height:        layout.Px(float64(height)),
			Padding:       layout.Spacing{Top: layout.Px(1), Right: layout.Px(2), Bottom: layout.Px(1), Left: layout.Px(2)},
		},
	}
	rootStyled := renderer.NewStyledNode(root, nil)

	// Header
	fgWhite, _ := color.ParseColor("#FFFFFF")
	bgBlue, _ := color.ParseColor("oklch(0.5 0.25 250)")
	borderBlue, _ := color.ParseColor("oklch(0.7 0.25 250)")

	headerNode := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   layout.Px(float64(width - 4)),
			Height:  layout.Px(3),
			Margin:  layout.Spacing{Bottom: layout.Px(1)},
		},
	}
	headerStyle := &renderer.Style{
		Foreground:  &fgWhite,
		Background:  &bgBlue,
		Bold:        true,
		BorderColor: &borderBlue,
		TextAlign:   renderer.TextAlignCenter,
	}
	headerStyle.WithBorder(renderer.ThickBorder)
	headerStyled := renderer.NewStyledNode(headerNode, headerStyle)
	headerStyled.Content = "System Dashboard"
	rootStyled.AddChild(headerStyled)

	// Gauges row
	gaugesContainer := &layout.Node{
		Style: layout.Style{
			Display:       layout.DisplayFlex,
			FlexDirection: layout.FlexDirectionRow,
			Width:         layout.Px(float64(width - 4)),
			Height:        layout.Px(6),
			Margin:        layout.Spacing{Bottom: layout.Px(1)},
		},
	}
	gaugesStyled := renderer.NewStyledNode(gaugesContainer, nil)

	// Create gauges
	gaugeWidth := (width - 4 - 6) / 4
	labels := []string{"CPU", "Memory", "Disk", "Network"}
	values := []float64{45.2, 68.5, 82.1, 23.7}

	for i, label := range labels {
		gauge := createSimpleGauge(label, values[i], gaugeWidth)
		if i < len(labels)-1 {
			gauge.Node.Style.Margin.Right = layout.Px(2)
		}
		gaugesStyled.AddChild(gauge)
	}

	rootStyled.AddChild(gaugesStyled)

	// Info panels - use FlexGrow to fill remaining space
	panelsContainer := &layout.Node{
		Style: layout.Style{
			Display:       layout.DisplayFlex,
			FlexDirection: layout.FlexDirectionRow,
			Width:         layout.Px(float64(width - 4)),
			FlexGrow:      1,
		},
	}
	panelsStyled := renderer.NewStyledNode(panelsContainer, nil)

	// Left panel
	leftPanel := createInfoPanel("System Info", "OS: darwin\nCPUs: 14\nMemory: 16GB", (width-4)/2)
	leftPanel.Node.Style.FlexGrow = 1
	leftPanel.Node.Style.Margin.Right = layout.Px(1)
	panelsStyled.AddChild(leftPanel)

	// Right panel
	rightPanel := createInfoPanel("Activity", "Service started\nConnection OK\nCache updated", (width-4)/2)
	rightPanel.Node.Style.FlexGrow = 1
	rightPanel.Node.Style.Margin.Left = layout.Px(1)
	panelsStyled.AddChild(rightPanel)

	rootStyled.AddChild(panelsStyled)

	// Footer
	footerNode := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   layout.Px(float64(width - 4)),
			Height:  layout.Px(1),
			Margin:  layout.Spacing{Top: layout.Px(1)},
		},
	}
	fgGray, _ := color.ParseColor("#888888")
	footerStyle := &renderer.Style{
		Foreground: &fgGray,
		Dim:        true,
		TextAlign:  renderer.TextAlignCenter,
	}
	footerStyled := renderer.NewStyledNode(footerNode, footerStyle)
	footerStyled.Content = "Dashboard Demo"
	rootStyled.AddChild(footerStyled)

	return rootStyled
}

func createSimpleGauge(label string, value float64, width int) *renderer.StyledNode {
	node := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   layout.Px(float64(width)),
			Height:  layout.Px(6),
		},
	}

	borderGray, _ := color.ParseColor("#5A5A5A")
	style := &renderer.Style{
		BorderColor: &borderGray,
	}
	style.WithBorder(renderer.NormalBorder)

	styledNode := renderer.NewStyledNode(node, style)

	// Build content
	content := fmt.Sprintf("%s\n%.1f%%\n", label, value)

	// Progress bar
	barWidth := width - 4
	filled := int(float64(barWidth) * value / 100.0)
	bar := ""
	for i := 0; i < barWidth; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}
	content += bar

	styledNode.Content = content

	return styledNode
}

func createInfoPanel(title, content string, width int) *renderer.StyledNode {
	node := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   layout.Px(float64(width)),
		},
	}

	fgWhite, _ := color.ParseColor("#FFFFFF")
	bgGreen, _ := color.ParseColor("oklch(0.4 0.15 150)")
	borderGreen, _ := color.ParseColor("oklch(0.6 0.15 150)")
	style := &renderer.Style{
		Foreground:  &fgWhite,
		Background:  &bgGreen,
		BorderColor: &borderGreen,
	}
	style.WithBorder(renderer.RoundedBorder)

	styledNode := renderer.NewStyledNode(node, style)
	styledNode.Content = title + "\n\n" + content

	return styledNode
}
