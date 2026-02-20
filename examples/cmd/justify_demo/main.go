package main

import (
	"fmt"

	"github.com/SCKelemen/color"
	"github.com/SCKelemen/layout"
	"github.com/SCKelemen/cli/renderer"
)

func main() {
	width, height := 80, 35

	screen := renderer.NewScreen(width, height)
	root := buildJustifyDemo(width, height)

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

func buildJustifyDemo(width, height int) *renderer.StyledNode {
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
	bgGreen, _ := color.ParseColor("oklch(0.5 0.2 150)")
	borderGreen, _ := color.ParseColor("oklch(0.7 0.2 150)")

	headerNode := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   layout.Px(float64(width - 4)),
			Height:  layout.Px(3),
			Margin:  layout.Spacing{Top: layout.Px(0), Right: layout.Px(0), Bottom: layout.Px(1), Left: layout.Px(0)},
		},
	}
	headerStyle := &renderer.Style{
		Foreground:  &fgWhite,
		Background:  &bgGreen,
		Bold:        true,
		BorderColor: &borderGreen,
		TextAlign:   renderer.TextAlignCenter,
	}
	headerStyle.WithBorder(renderer.RoundedBorder)
	headerStyled := renderer.NewStyledNode(headerNode, headerStyle)
	headerStyled.Content = "Text Justification Demo"
	rootStyled.AddChild(headerStyled)

	borderGray, _ := color.ParseColor("#5A5A5A")

	// Multi-paragraph text for justify demo
	paragraphText := `The quick brown fox jumps over the lazy dog. This is a longer piece of text that will demonstrate text justification with proper word spacing.

Typography has a long history of justification used in printed books and newspapers. The goal is to create clean right margins by distributing space between words.

Notice how the last line of each paragraph is left-aligned, following typographic convention. Only the full lines are justified.`

	// Left-aligned (for comparison)
	leftNode := createTextBox(width-4, 7, "Left Aligned (No Justification)", paragraphText,
		renderer.TextAlignLeft, renderer.TextWrapNormal, &borderGray)
	rootStyled.AddChild(leftNode)

	// Justified
	justifyNode := createTextBox(width-4, 7, "Justified (Space Between Words)", paragraphText,
		renderer.TextAlignJustify, renderer.TextWrapNormal, &borderGray)
	rootStyled.AddChild(justifyNode)

	// Footer
	footerNode := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   layout.Px(float64(width - 4)),
			Height:  layout.Px(2),
			Margin:  layout.Spacing{Top: layout.Px(1), Right: layout.Px(0), Bottom: layout.Px(0), Left: layout.Px(0)},
		},
	}
	fgGray, _ := color.ParseColor("#888888")
	footerStyle := &renderer.Style{
		Foreground: &fgGray,
		Dim:        true,
		TextAlign:  renderer.TextAlignCenter,
	}
	footerStyled := renderer.NewStyledNode(footerNode, footerStyle)
	footerStyled.Content = "Justification adds space between words to align both left and right edges"
	rootStyled.AddChild(footerStyled)

	return rootStyled
}

func createTextBox(width, height int, title, content string, align renderer.TextAlign, wrap renderer.TextWrap, borderColor *color.Color) *renderer.StyledNode {
	node := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   layout.Px(float64(width)),
			Height:  layout.Px(float64(height)),
			Margin:  layout.Spacing{Top: layout.Px(1), Right: layout.Px(0), Bottom: layout.Px(0), Left: layout.Px(0)},
		},
	}

	style := &renderer.Style{
		BorderColor: borderColor,
		TextAlign:   align,
		TextWrap:    wrap,
	}
	style.WithBorder(renderer.NormalBorder)

	styledNode := renderer.NewStyledNode(node, style)
	styledNode.Content = title + "\n" + content

	return styledNode
}
