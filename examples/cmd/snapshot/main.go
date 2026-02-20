package main

import (
	"fmt"

	"github.com/SCKelemen/cli/components"
	"github.com/SCKelemen/cli/renderer"
	"github.com/SCKelemen/color"
	"github.com/SCKelemen/layout"
)

func main() {
	// Terminal dimensions
	width, height := 80, 30

	// Create screen
	screen := renderer.NewScreen(width, height)

	// Build layout
	root := buildLayout(width, height)

	// Perform layout calculation
	constraints := layout.Tight(float64(width), float64(height))
	ctx := &layout.LayoutContext{
		ViewportWidth:  float64(width),
		ViewportHeight: float64(height),
		RootFontSize:   16,
	}
	layout.Layout(root.Node, constraints, ctx)

	// Render to screen
	screen.Render(root)

	// Output
	fmt.Print(screen.String())
}

func buildLayout(width, height int) *renderer.StyledNode {
	// Create root container with flexbox column layout
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
			Margin:  layout.Spacing{Top: layout.Px(0), Right: layout.Px(0), Bottom: layout.Px(1), Left: layout.Px(0)},
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
		"Press '1' or '2' to toggle sections\nPress 'q' to quit",
	).WithBorderColor(&borderGray)
	rootStyled.AddChild(info.ToStyledNode())

	// Loading indicators container
	loadingContainer := &layout.Node{
		Style: layout.Style{
			Display:       layout.DisplayFlex,
			FlexDirection: layout.FlexDirectionRow,
			Margin:        layout.Spacing{Top: layout.Px(1), Right: layout.Px(0), Bottom: layout.Px(1), Left: layout.Px(0)},
		},
	}
	loadingStyled := renderer.NewStyledNode(loadingContainer, nil)

	// Add spinner
	spinner := components.NewSpinnerDots()
	spinnerNode := spinner.ToStyledNode()
	spinnerNode.Node.Style.Margin = layout.Spacing{Top: layout.Px(0), Right: layout.Px(2), Bottom: layout.Px(0), Left: layout.Px(0)}
	loadingStyled.AddChild(spinnerNode)

	// Add loading dots
	loading := components.NewLoadingDots()
	loadingNode := loading.ToStyledNode()
	loadingNode.Node.Style.Margin = layout.Spacing{Top: layout.Px(0), Right: layout.Px(2), Bottom: layout.Px(0), Left: layout.Px(0)}
	loadingStyled.AddChild(loadingNode)

	rootStyled.AddChild(loadingStyled)

	// Progress bar
	progress := components.NewProgressBar(40)
	progress.SetProgress(0.65)
	progressNode := progress.ToStyledNode()
	progressNode.Node.Style.Margin = layout.Spacing{Top: layout.Px(1), Right: layout.Px(0), Bottom: layout.Px(1), Left: layout.Px(0)}
	rootStyled.AddChild(progressNode)

	// Collapsible sections
	section1 := components.NewCollapsible(
		"Project Information",
		"This is a proof of concept for integrating\nthe layout engine with terminal UIs.\n\nWe're building beautiful terminal interfaces!",
	)
	section1Node := section1.ToStyledNode()
	section1Node.Node.Style.Margin = layout.Spacing{Top: layout.Px(1), Right: layout.Px(0), Bottom: layout.Px(0), Left: layout.Px(0)}
	rootStyled.AddChild(section1Node)

	section2 := components.NewCollapsible(
		"Features",
		"- Responsive layouts using CSS Grid/Flexbox\n- OKLCH color gradients\n- Collapsible sections\n- Animated loading indicators\n- Dynamic resizing",
	)
	section2.Expanded = false // Collapse this one for demo
	section2Node := section2.ToStyledNode()
	section2Node.Node.Style.Margin = layout.Spacing{Top: layout.Px(1), Right: layout.Px(0), Bottom: layout.Px(0), Left: layout.Px(0)}
	rootStyled.AddChild(section2Node)

	// Footer
	footerNode := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   layout.Px(float64(width - 4)),
			Height:  layout.Px(3),
			Margin:  layout.Spacing{Top: layout.Px(1), Right: layout.Px(0), Bottom: layout.Px(0), Left: layout.Px(0)},
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
