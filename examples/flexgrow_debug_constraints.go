package main

import (
	"fmt"

	"github.com/SCKelemen/layout"
)

func main() {
	width, height := 100, 30

	fmt.Println("=== Debugging Nested FlexGrow Constraints ===\n")

	// Create a simple nested structure similar to dashboard
	root := &layout.Node{
		Style: layout.Style{
			Display:       layout.DisplayFlex,
			FlexDirection: layout.FlexDirectionColumn,
			Width:         layout.Px(float64(width)),
			Height:        layout.Px(float64(height)),
		},
	}

	// Fixed header
	header := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   layout.Px(float64(width)),
			Height:  layout.Px(5),
		},
	}
	root.AddChild(header)

	// Row container with FlexGrow
	rowContainer := &layout.Node{
		Style: layout.Style{
			Display:       layout.DisplayFlex,
			FlexDirection: layout.FlexDirectionRow,
			Width:         layout.Px(float64(width)),
			FlexGrow:      1, // Should fill remaining space (25 lines)
		},
	}
	root.AddChild(rowContainer)

	// Left panel with FlexGrow but no explicit dimensions
	leftPanel := &layout.Node{
		Style: layout.Style{
			Display:  layout.DisplayBlock,
			FlexGrow: 1,
			// No Width, no Height
		},
	}
	rowContainer.AddChild(leftPanel)

	// Right panel with FlexGrow but no explicit dimensions
	rightPanel := &layout.Node{
		Style: layout.Style{
			Display:  layout.DisplayBlock,
			FlexGrow: 1,
			// No Width, no Height
		},
	}
	rowContainer.AddChild(rightPanel)

	// Fixed footer
	footer := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   layout.Px(float64(width)),
			Height:  layout.Px(3),
		},
	}
	root.AddChild(footer)

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
	fmt.Printf("Header: %v\n", header.Rect)
	fmt.Printf("Row Container: %v\n", rowContainer.Rect)
	fmt.Printf("  -> Expected: {0 5 100 22} (should fill remaining height)\n")
	fmt.Printf("Left Panel: %v\n", leftPanel.Rect)
	fmt.Printf("  -> Expected: {0 0 50 22} (should inherit parent height and stretch)\n")
	fmt.Printf("Right Panel: %v\n", rightPanel.Rect)
	fmt.Printf("  -> Expected: {50 0 50 22} (should inherit parent height and stretch)\n")
	fmt.Printf("Footer: %v\n", footer.Rect)

	// Diagnosis
	fmt.Println("\n=== Diagnosis ===")
	if rowContainer.Rect.Height < 20 {
		fmt.Println("✗ Row container didn't grow to fill available space")
	} else {
		fmt.Println("✓ Row container grew correctly to", rowContainer.Rect.Height)
	}

	if leftPanel.Rect.Height == 0 {
		fmt.Println("✗ Left panel has zero height - AlignItems:Stretch not working!")
		fmt.Println("  This is the bug: children should inherit parent's computed height")
	} else if leftPanel.Rect.Height != rowContainer.Rect.Height {
		fmt.Printf("✗ Left panel height (%v) doesn't match container (%v)\n",
			leftPanel.Rect.Height, rowContainer.Rect.Height)
	} else {
		fmt.Println("✓ Left panel stretched to match container")
	}

	if rightPanel.Rect.Height == 0 {
		fmt.Println("✗ Right panel has zero height - AlignItems:Stretch not working!")
	} else if rightPanel.Rect.Height != rowContainer.Rect.Height {
		fmt.Printf("✗ Right panel height (%v) doesn't match container (%v)\n",
			rightPanel.Rect.Height, rowContainer.Rect.Height)
	} else {
		fmt.Println("✓ Right panel stretched to match container")
	}
}
