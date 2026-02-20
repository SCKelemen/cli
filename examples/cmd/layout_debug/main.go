package main

import (
	"fmt"

	"github.com/SCKelemen/color"
	"github.com/SCKelemen/layout"
	"github.com/SCKelemen/cli/renderer"
)

func main() {
	width, height := 100, 30

	// Test 1: FlexGrow without MinHeight
	fmt.Println("=== Test 1: FlexGrow without MinHeight ===")
	root1 := &layout.Node{
		Style: layout.Style{
			Display:       layout.DisplayFlex,
			FlexDirection: layout.FlexDirectionColumn,
			Width:         layout.Px(float64(width)),
			Height:        layout.Px(float64(height)),
		},
	}

	header1 := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   layout.Px(float64(width)),
			Height:  layout.Px(5),
		},
	}
	root1.AddChild(header1)

	content1 := &layout.Node{
		Style: layout.Style{
			Display:  layout.DisplayBlock,
			Width:    layout.Px(float64(width)),
			FlexGrow: 1,
		},
	}
	root1.AddChild(content1)

	footer1 := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   layout.Px(float64(width)),
			Height:  layout.Px(3),
		},
	}
	root1.AddChild(footer1)

	constraints := layout.Tight(float64(width), float64(height))
	ctx := &layout.LayoutContext{
		ViewportWidth:  float64(width),
		ViewportHeight: float64(height),
		RootFontSize:   16,
	}
	layout.Layout(root1, constraints, ctx)

	fmt.Printf("Root: %v\n", root1.Rect)
	fmt.Printf("Header: %v (expected height: 5)\n", header1.Rect)
	fmt.Printf("Content: %v (expected height: ~22, to fill remaining)\n", content1.Rect)
	fmt.Printf("Footer: %v (expected height: 3)\n", footer1.Rect)

	// Test 2: FlexGrow with MinHeight
	fmt.Println("\n=== Test 2: FlexGrow with MinHeight ===")
	root2 := &layout.Node{
		Style: layout.Style{
			Display:       layout.DisplayFlex,
			FlexDirection: layout.FlexDirectionColumn,
			Width:         layout.Px(float64(width)),
			Height:        layout.Px(float64(height)),
		},
	}

	header2 := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   layout.Px(float64(width)),
			Height:  layout.Px(5),
		},
	}
	root2.AddChild(header2)

	content2 := &layout.Node{
		Style: layout.Style{
			Display:   layout.DisplayBlock,
			Width:     layout.Px(float64(width)),
			FlexGrow:  1,
			MinHeight: layout.Px(10),
		},
	}
	root2.AddChild(content2)

	footer2 := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   layout.Px(float64(width)),
			Height:  layout.Px(3),
		},
	}
	root2.AddChild(footer2)

	layout.Layout(root2, constraints, ctx)

	fmt.Printf("Root: %v\n", root2.Rect)
	fmt.Printf("Header: %v (expected height: 5)\n", header2.Rect)
	fmt.Printf("Content: %v (expected height: ~22, to fill remaining)\n", content2.Rect)
	fmt.Printf("Footer: %v (expected height: 3)\n", footer2.Rect)

	// Test 3: Nested FlexGrow - row with two FlexGrow children
	fmt.Println("\n=== Test 3: Nested FlexGrow (Row with FlexGrow children) ===")
	root3 := &layout.Node{
		Style: layout.Style{
			Display:       layout.DisplayFlex,
			FlexDirection: layout.FlexDirectionColumn,
			Width:         layout.Px(float64(width)),
			Height:        layout.Px(float64(height)),
		},
	}

	row := &layout.Node{
		Style: layout.Style{
			Display:       layout.DisplayFlex,
			FlexDirection: layout.FlexDirectionRow,
			Width:         layout.Px(float64(width)),
			FlexGrow:      1,
		},
	}
	root3.AddChild(row)

	leftPanel := &layout.Node{
		Style: layout.Style{
			Display:  layout.DisplayBlock,
			Width:    layout.Px(float64(width / 2)),
			FlexGrow: 1,
		},
	}
	row.AddChild(leftPanel)

	rightPanel := &layout.Node{
		Style: layout.Style{
			Display:  layout.DisplayBlock,
			Width:    layout.Px(float64(width / 2)),
			FlexGrow: 1,
		},
	}
	row.AddChild(rightPanel)

	layout.Layout(root3, constraints, ctx)

	fmt.Printf("Root: %v\n", root3.Rect)
	fmt.Printf("Row: %v (expected height: 30, to fill all space)\n", row.Rect)
	fmt.Printf("Left Panel: %v (expected height: 30)\n", leftPanel.Rect)
	fmt.Printf("Right Panel: %v (expected height: 30)\n", rightPanel.Rect)

	// Now render to see what actually displays
	fmt.Println("\n=== Rendering Test 2 (with MinHeight) ===")
	screen := renderer.NewScreen(width, height)

	rootStyled := renderer.NewStyledNode(root2, nil)

	borderGray, _ := color.ParseColor("#5A5A5A")
	headerStyled := renderer.NewStyledNode(header2, &renderer.Style{BorderColor: &borderGray})
	headerStyled.Style.WithBorder(renderer.NormalBorder)
	headerStyled.Content = "HEADER"
	rootStyled.AddChild(headerStyled)

	bgGreen, _ := color.ParseColor("oklch(0.4 0.15 150)")
	contentStyled := renderer.NewStyledNode(content2, &renderer.Style{Background: &bgGreen, BorderColor: &borderGray})
	contentStyled.Style.WithBorder(renderer.RoundedBorder)
	contentStyled.Content = "CONTENT AREA\nShould fill space"
	rootStyled.AddChild(contentStyled)

	footerStyled := renderer.NewStyledNode(footer2, &renderer.Style{BorderColor: &borderGray})
	footerStyled.Style.WithBorder(renderer.NormalBorder)
	footerStyled.Content = "FOOTER"
	rootStyled.AddChild(footerStyled)

	screen.Render(rootStyled)
	fmt.Print(screen.String())
}
