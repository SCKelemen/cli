package main

import (
	"fmt"
	"time"

	"github.com/SCKelemen/color"
	"github.com/SCKelemen/layout"
	"github.com/SCKelemen/cli/components"
	"github.com/SCKelemen/cli/renderer"
)

func main() {
	// Terminal dimensions
	width, height := 80, 35

	// Create components
	loading := components.NewLoadingDots()
	spinner := components.NewSpinnerDots()
	progress := components.NewProgressBar(50)

	section1 := components.NewCollapsible(
		"Project Information",
		"This is a proof of concept for integrating\nthe layout engine with terminal UIs.\n\nWe're building beautiful terminal interfaces!",
	)

	section2 := components.NewCollapsible(
		"Features",
		"- Responsive layouts using CSS Grid/Flexbox\n- OKLCH color gradients\n- Collapsible sections\n- Animated loading indicators\n- Dynamic resizing",
	)

	// Render 10 frames to show animation
	fmt.Println("=== Animation Test: Showing 10 frames ===\n")

	for frame := 0; frame < 10; frame++ {
		fmt.Printf("--- Frame %d ---\n", frame+1)

		// Update animations
		now := time.Now()
		loading.Update(now)
		spinner.Update(now)
		progress.SetProgress(float64(frame) / 9.0)

		// Toggle sections on certain frames
		if frame == 3 {
			section1.Toggle() // Collapse
		}
		if frame == 6 {
			section1.Toggle() // Expand
			section2.Toggle() // Expand
		}

		// Build layout
		root := buildFullLayout(width, height, loading, spinner, progress, section1, section2)

		// Create screen and render
		screen := renderer.NewScreen(width, height)
		constraints := layout.Tight(float64(width), float64(height))
		ctx := &layout.LayoutContext{
			ViewportWidth:  float64(width),
			ViewportHeight: float64(height),
			RootFontSize:   16,
		}
		layout.Layout(root.Node, constraints, ctx)
		screen.Render(root)

		// Output
		fmt.Print(screen.String())
		fmt.Println()

		// Advance time for next frame
		time.Sleep(100 * time.Millisecond)
	}
}

func buildFullLayout(width, height int, loading *components.LoadingDots, spinner *components.SpinnerDots, progress *components.ProgressBar, section1, section2 *components.Collapsible) *renderer.StyledNode {
	// Create root container
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
	headerNode := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   layout.Px(float64(width - 4)),
			Height:  layout.Px(5),
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
	headerStyled.Content = "  Terminal Layout Engine Demo  "
	rootStyled.AddChild(headerStyled)

	// Info message
	borderGray, _ := color.ParseColor("#5A5A5A")
	info := components.NewMessageBlock(
		"Animation Test - Loading indicators and collapsible sections",
	).WithBorderColor(&borderGray)
	rootStyled.AddChild(info.ToStyledNode())

	// Loading indicators container
	loadingContainer := &layout.Node{
		Style: layout.Style{
			Display:       layout.DisplayFlex,
			FlexDirection: layout.FlexDirectionRow,
		},
	}
	loadingStyled := renderer.NewStyledNode(loadingContainer, nil)

	// Add spinner
	spinnerNode := spinner.ToStyledNode()
	spinnerNode.Node.Style.Margin = layout.Spacing{Top: layout.Px(0), Right: layout.Px(2), Bottom: layout.Px(0), Left: layout.Px(0)}
	loadingStyled.AddChild(spinnerNode)

	// Add loading dots
	loadingNode := loading.ToStyledNode()
	loadingNode.Node.Style.Margin = layout.Spacing{Top: layout.Px(0), Right: layout.Px(2), Bottom: layout.Px(0), Left: layout.Px(0)}
	loadingStyled.AddChild(loadingNode)

	rootStyled.AddChild(loadingStyled)

	// Progress bar
	progressNode := progress.ToStyledNode()
	rootStyled.AddChild(progressNode)

	// Collapsible sections
	section1Node := section1.ToStyledNode()
	rootStyled.AddChild(section1Node)

	section2Node := section2.ToStyledNode()
	rootStyled.AddChild(section2Node)

	// Footer
	footerNode := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   layout.Px(float64(width - 4)),
			Height:  layout.Px(1),
		},
	}
	fgGray, _ := color.ParseColor("#888888")
	footerStyle := &renderer.Style{
		Foreground: &fgGray,
		Dim:        true,
	}
	footerStyled := renderer.NewStyledNode(footerNode, footerStyle)
	footerStyled.Content = fmt.Sprintf("Terminal: %dx%d", width, height)
	rootStyled.AddChild(footerStyled)

	return rootStyled
}
