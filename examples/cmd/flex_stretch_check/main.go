package main

import (
	"fmt"

	"github.com/SCKelemen/cli/renderer"
	"github.com/SCKelemen/color"
	"github.com/SCKelemen/layout"
)

func main() {
	width, height := 100, 30

	fmt.Print("=== Test: AlignItems with FlexGrow container ===\n")

	// Test 1: With explicit AlignItems: Stretch
	root1 := &layout.Node{
		Style: layout.Style{
			Display:       layout.DisplayFlex,
			FlexDirection: layout.FlexDirectionColumn,
			Width:         layout.Px(float64(width)),
			Height:        layout.Px(float64(height)),
		},
	}
	rootStyled1 := renderer.NewStyledNode(root1, nil)

	header1 := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   layout.Px(float64(width)),
			Height:  layout.Px(5),
		},
	}
	headerStyled1 := renderer.NewStyledNode(header1, nil)
	headerStyled1.Content = "HEADER"
	rootStyled1.AddChild(headerStyled1)

	// Row container with FlexGrow AND explicit AlignItems
	row1 := &layout.Node{
		Style: layout.Style{
			Display:       layout.DisplayFlex,
			FlexDirection: layout.FlexDirectionRow,
			Width:         layout.Px(float64(width)),
			FlexGrow:      1,
			AlignItems:    layout.AlignItemsStretch, // Explicit stretch
		},
	}
	rowStyled1 := renderer.NewStyledNode(row1, nil)

	left1 := &layout.Node{
		Style: layout.Style{
			Display:  layout.DisplayBlock,
			FlexGrow: 1,
		},
	}
	bgGreen, _ := color.ParseColor("#00FF00")
	leftStyled1 := renderer.NewStyledNode(left1, &renderer.Style{Background: &bgGreen})
	leftStyled1.Content = "LEFT (stretch)"
	rowStyled1.AddChild(leftStyled1)

	right1 := &layout.Node{
		Style: layout.Style{
			Display:  layout.DisplayBlock,
			FlexGrow: 1,
		},
	}
	bgBlue, _ := color.ParseColor("#0000FF")
	rightStyled1 := renderer.NewStyledNode(right1, &renderer.Style{Background: &bgBlue})
	rightStyled1.Content = "RIGHT (stretch)"
	rowStyled1.AddChild(rightStyled1)

	rootStyled1.AddChild(rowStyled1)

	footer1 := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   layout.Px(float64(width)),
			Height:  layout.Px(3),
		},
	}
	footerStyled1 := renderer.NewStyledNode(footer1, nil)
	footerStyled1.Content = "FOOTER"
	rootStyled1.AddChild(footerStyled1)

	constraints := layout.Tight(float64(width), float64(height))
	ctx := &layout.LayoutContext{
		ViewportWidth:  float64(width),
		ViewportHeight: float64(height),
		RootFontSize:   16,
	}
	layout.Layout(root1, constraints, ctx)

	fmt.Println("Test 1: Explicit AlignItems:Stretch on FlexGrow container")
	fmt.Printf("Row: %v\n", row1.Rect)
	fmt.Printf("Left: %v\n", left1.Rect)
	fmt.Printf("Right: %v\n", right1.Rect)

	if left1.Rect.Height > 15 {
		fmt.Println("✓ LEFT PANEL STRETCHED!")
	} else {
		fmt.Printf("✗ Left panel height is only %v\n", left1.Rect.Height)
	}

	fmt.Print("\n=== Rendering ===")
	screen := renderer.NewScreen(width, height)
	screen.Render(rootStyled1)
	fmt.Print(screen.String())
}
