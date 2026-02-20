package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/SCKelemen/cli/components"
	"github.com/SCKelemen/cli/renderer"
	"github.com/SCKelemen/color"
	"github.com/SCKelemen/layout"
)

func main() {
	fmt.Print("=== Terminal Layout Engine - Animation Demo ===\n")
	fmt.Print("Demonstrating animated components over 8 frames...\n")

	// Create animated components
	loading := components.NewLoadingDots()
	spinner := components.NewSpinnerDots()
	progress := components.NewProgressBar(60)

	for frame := 0; frame < 8; frame++ {
		fmt.Printf("Frame %d:\n", frame+1)
		fmt.Println(strings.Repeat("─", 80))

		// Update animations
		now := time.Now().Add(time.Duration(frame*150) * time.Millisecond)
		loading.Update(now)
		spinner.Update(now)
		progress.SetProgress(float64(frame) / 7.0)

		// Render frame
		screen := renderFrame(loading, spinner, progress, frame)
		fmt.Println(screen)
		fmt.Println()

		time.Sleep(150 * time.Millisecond)
	}

	fmt.Println("Animation complete!")
}

func renderFrame(loading *components.LoadingDots, spinner *components.SpinnerDots, progress *components.ProgressBar, frame int) string {
	width, height := 80, 20

	// Create layout
	root := &layout.Node{
		Style: layout.Style{
			Display:       layout.DisplayFlex,
			FlexDirection: layout.FlexDirectionColumn,
			Width:         layout.Px(float64(width)),
			Height:        layout.Px(float64(height)),
			Padding:       layout.Uniform(layout.Px(2)),
		},
	}
	rootStyled := renderer.NewStyledNode(root, nil)

	// Title
	titleNode := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   layout.Px(float64(width - 4)),
			Height:  layout.Px(5),
		},
	}
	fgWhite, _ := color.ParseColor("#FFFFFF")
	bgPurple, _ := color.ParseColor("oklch(0.5 0.2 270)")
	borderPurple, _ := color.ParseColor("oklch(0.7 0.2 270)")
	titleStyle := &renderer.Style{
		Foreground:  &fgWhite,
		Background:  &bgPurple,
		Bold:        true,
		BorderColor: &borderPurple,
	}
	titleStyle.WithBorder(renderer.DoubleBorder)
	titleStyled := renderer.NewStyledNode(titleNode, titleStyle)
	titleStyled.Content = fmt.Sprintf("  Animation Frame %d/8  ", frame+1)
	rootStyled.AddChild(titleStyled)

	// Spinner
	spinnerWrapper := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   layout.Px(float64(width - 4)),
			Height:  layout.Px(1),
		},
	}
	spinnerWrapperStyled := renderer.NewStyledNode(spinnerWrapper, nil)

	spinnerNode := spinner.ToStyledNode()
	spinnerNode.Node.Style.Width = layout.Px(2)
	spinnerWrapperStyled.AddChild(spinnerNode)

	// Add label
	labelNode := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   layout.Px(20),
			Height:  layout.Px(1),
		},
	}
	fgLabel, _ := color.ParseColor("#AAAAAA")
	labelStyle := &renderer.Style{
		Foreground: &fgLabel,
	}
	labelStyled := renderer.NewStyledNode(labelNode, labelStyle)
	labelStyled.Content = " Spinner animation"
	spinnerWrapperStyled.AddChild(labelStyled)

	rootStyled.AddChild(spinnerWrapperStyled)

	// Loading dots
	loadingWrapper := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   layout.Px(float64(width - 4)),
			Height:  layout.Px(1),
		},
	}
	loadingWrapperStyled := renderer.NewStyledNode(loadingWrapper, nil)

	loadingNode := loading.ToStyledNode()
	loadingWrapperStyled.AddChild(loadingNode)

	rootStyled.AddChild(loadingWrapperStyled)

	// Progress bar
	progressWrapper := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   layout.Px(float64(width - 4)),
			Height:  layout.Px(3),
		},
	}
	progressWrapperStyled := renderer.NewStyledNode(progressWrapper, nil)

	progressLabelNode := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   layout.Px(float64(width - 4)),
			Height:  layout.Px(1),
		},
	}
	fgProgress, _ := color.ParseColor("#AAAAAA")
	progressLabelStyle := &renderer.Style{
		Foreground: &fgProgress,
	}
	progressLabelStyled := renderer.NewStyledNode(progressLabelNode, progressLabelStyle)
	progressLabelStyled.Content = fmt.Sprintf("Progress: %.0f%%", progress.Progress*100)
	progressWrapperStyled.AddChild(progressLabelStyled)

	progressNode := progress.ToStyledNode()
	progressWrapperStyled.AddChild(progressNode)

	rootStyled.AddChild(progressWrapperStyled)

	// Footer
	footerNode := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   layout.Px(float64(width - 4)),
			Height:  layout.Px(1),
		},
	}
	fgGray, _ := color.ParseColor("#666666")
	footerStyle := &renderer.Style{
		Foreground: &fgGray,
		Dim:        true,
	}
	footerStyled := renderer.NewStyledNode(footerNode, footerStyle)
	footerStyled.Content = "Using OKLCH colors, CSS layout, and 30fps animations"
	rootStyled.AddChild(footerStyled)

	// Perform layout
	constraints := layout.Tight(float64(width), float64(height))
	ctx := &layout.LayoutContext{
		ViewportWidth:  float64(width),
		ViewportHeight: float64(height),
		RootFontSize:   16,
	}
	layout.Layout(root, constraints, ctx)

	// Render
	screen := renderer.NewScreen(width, height)
	screen.Render(rootStyled)

	return screen.String()
}
