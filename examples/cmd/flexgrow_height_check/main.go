package main

import (
	"fmt"

	"github.com/SCKelemen/cli/renderer"
	"github.com/SCKelemen/color"
	"github.com/SCKelemen/layout"
)

func main() {
	width, height := 100, 30

	fmt.Println("=== Test: Row container with FlexGrow children ===")
	fmt.Println("Container should have height 20")
	fmt.Println("Children with FlexGrow should inherit that height")
	fmt.Println()

	root := &layout.Node{
		Style: layout.Style{
			Display:       layout.DisplayFlex,
			FlexDirection: layout.FlexDirectionColumn,
			Width:         layout.Px(float64(width)),
			Height:        layout.Px(float64(height)),
		},
	}
	rootStyled := renderer.NewStyledNode(root, nil)

	// Fixed header
	header := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   layout.Px(float64(width)),
			Height:  layout.Px(5),
		},
	}
	headerStyled := renderer.NewStyledNode(header, nil)
	headerStyled.Content = "HEADER"
	rootStyled.AddChild(headerStyled)

	// Row container with explicit height
	rowContainer := &layout.Node{
		Style: layout.Style{
			Display:       layout.DisplayFlex,
			FlexDirection: layout.FlexDirectionRow,
			Width:         layout.Px(float64(width)),
			Height:        layout.Px(20), // Explicit height
		},
	}
	rowStyled := renderer.NewStyledNode(rowContainer, nil)

	// Left panel - NO MinHeight, only FlexGrow
	leftPanel := &layout.Node{
		Style: layout.Style{
			Display:  layout.DisplayBlock,
			FlexGrow: 1, // Should grow in width
			// NO Width specified
			// NO Height specified
			// NO MinHeight specified
		},
	}
	bgGreen, _ := color.ParseColor("#00FF00")
	leftStyled := renderer.NewStyledNode(leftPanel, &renderer.Style{
		Background: &bgGreen,
	})
	leftStyled.Content = "LEFT PANEL\nShould be 20 lines tall"
	rowStyled.AddChild(leftStyled)

	// Right panel - NO MinHeight, only FlexGrow
	rightPanel := &layout.Node{
		Style: layout.Style{
			Display:  layout.DisplayBlock,
			FlexGrow: 1, // Should grow in width
			// NO Width specified
			// NO Height specified
			// NO MinHeight specified
		},
	}
	bgBlue, _ := color.ParseColor("#0000FF")
	rightStyled := renderer.NewStyledNode(rightPanel, &renderer.Style{
		Background: &bgBlue,
	})
	rightStyled.Content = "RIGHT PANEL\nShould be 20 lines tall"
	rowStyled.AddChild(rightStyled)

	rootStyled.AddChild(rowStyled)

	// Fixed footer
	footer := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   layout.Px(float64(width)),
			Height:  layout.Px(3),
		},
	}
	footerStyled := renderer.NewStyledNode(footer, nil)
	footerStyled.Content = "FOOTER"
	rootStyled.AddChild(footerStyled)

	// Layout
	constraints := layout.Tight(float64(width), float64(height))
	ctx := &layout.LayoutContext{
		ViewportWidth:  float64(width),
		ViewportHeight: float64(height),
		RootFontSize:   16,
	}
	layout.Layout(root, constraints, ctx)

	// Check results
	fmt.Printf("Root: %v\n", root.Rect)
	fmt.Printf("Header: %v (expected height: 5)\n", header.Rect)
	fmt.Printf("Row Container: %v (expected height: 20)\n", rowContainer.Rect)
	fmt.Printf("Left Panel: %v (expected height: 20 to match parent)\n", leftPanel.Rect)
	fmt.Printf("Right Panel: %v (expected height: 20 to match parent)\n", rightPanel.Rect)
	fmt.Printf("Footer: %v (expected height: 3)\n", footer.Rect)

	// Render to see what it looks like
	fmt.Print("\n=== Rendered Output ===")
	screen := renderer.NewScreen(width, height)
	screen.Render(rootStyled)
	fmt.Print(screen.String())
}
