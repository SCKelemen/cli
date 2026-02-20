package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/SCKelemen/cli/renderer"
	"github.com/SCKelemen/color"
	"github.com/SCKelemen/layout"
)

func main() {
	width, height := 100, 30
	rand.Seed(time.Now().UnixNano())

	root := buildDashboard(width, height)

	constraints := layout.Tight(float64(width), float64(height))
	ctx := &layout.LayoutContext{
		ViewportWidth:  float64(width),
		ViewportHeight: float64(height),
		RootFontSize:   16,
	}
	layout.Layout(root.Node, constraints, ctx)

	fmt.Printf("Terminal: %dx%d\n\n", width, height)
	fmt.Printf("Root: %v\n", root.Node.Rect)

	for i, child := range root.Children {
		fmt.Printf("\nChild %d: %v\n", i, child.Node.Rect)
		if len(child.Children) > 0 {
			for j, grandchild := range child.Children {
				fmt.Printf("  Grandchild %d: %v\n", j, grandchild.Node.Rect)
			}
		}
	}

	// Calculate what the heights should be
	fmt.Print("\n=== Expected Heights ===")
	fmt.Println("Root: 100x30")
	fmt.Println("With padding (1,2,1,2), content area: 96x28")
	fmt.Println("Header: 96x3 + margin 1 = 4 lines")
	fmt.Println("Gauges: 96x6 + margin 1 = 7 lines")
	fmt.Println("Footer: 96x1 + margin 1 = 2 lines")
	fmt.Println("Total fixed: 4 + 7 + 2 = 13 lines")
	fmt.Println("Remaining for panels with FlexGrow: 28 - 13 = 15 lines")
	fmt.Println("Each panel should get: 15 lines height")
}

func buildDashboard(width, height int) *renderer.StyledNode {
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
	leftPanel := createInfoPanel("System Info", "OS: darwin\nCPUs: 14", (width-4)/2)
	leftPanel.Node.Style.FlexGrow = 1
	leftPanel.Node.Style.Margin.Right = layout.Px(1)
	panelsStyled.AddChild(leftPanel)

	// Right panel
	rightPanel := createInfoPanel("Activity", "Service started", (width-4)/2)
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

func createInfoPanel(title, content string, width int) *renderer.StyledNode {
	node := &layout.Node{
		Style: layout.Style{
			Display:   layout.DisplayBlock,
			Width:     layout.Px(float64(width)),
			MinHeight: layout.Px(10),
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
