package main

import (
	"fmt"

	"github.com/SCKelemen/cli/renderer"
	"github.com/SCKelemen/color"
	"github.com/SCKelemen/layout"
)

func main() {
	width, height := 80, 45

	screen := renderer.NewScreen(width, height)
	root := buildWrapAlignDemo(width, height)

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

func buildWrapAlignDemo(width, height int) *renderer.StyledNode {
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
	bgPurple, _ := color.ParseColor("oklch(0.5 0.2 270)")
	borderPurple, _ := color.ParseColor("oklch(0.7 0.2 270)")

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
		Background:  &bgPurple,
		Bold:        true,
		BorderColor: &borderPurple,
		TextAlign:   renderer.TextAlignCenter,
	}
	headerStyle.WithBorder(renderer.RoundedBorder)
	headerStyled := renderer.NewStyledNode(headerNode, headerStyle)
	headerStyled.Content = "Text Wrapping & Alignment Demo"
	rootStyled.AddChild(headerStyled)

	borderGray, _ := color.ParseColor("#5A5A5A")
	longText := "The quick brown fox jumps over the lazy dog. This is a longer piece of text that will demonstrate various wrapping and alignment modes."

	// Left-aligned with normal wrapping
	leftNode := createTextBox(width-4, 5, "Left Aligned (Normal Wrap)", longText,
		renderer.TextAlignLeft, renderer.TextWrapNormal, &borderGray)
	rootStyled.AddChild(leftNode)

	// Center-aligned with normal wrapping
	centerNode := createTextBox(width-4, 5, "Center Aligned (Normal Wrap)", longText,
		renderer.TextAlignCenter, renderer.TextWrapNormal, &borderGray)
	rootStyled.AddChild(centerNode)

	// Right-aligned with normal wrapping
	rightNode := createTextBox(width-4, 5, "Right Aligned (Normal Wrap)", longText,
		renderer.TextAlignRight, renderer.TextWrapNormal, &borderGray)
	rootStyled.AddChild(rightNode)

	// Pretty wrapping (Knuth-Plass algorithm)
	prettyNode := createTextBox(width-4, 5, "Left Aligned (Pretty Wrap - Knuth-Plass)", longText,
		renderer.TextAlignLeft, renderer.TextWrapPretty, &borderGray)
	rootStyled.AddChild(prettyNode)

	// Center with pretty wrapping
	centerPrettyNode := createTextBox(width-4, 5, "Center Aligned (Pretty Wrap)", longText,
		renderer.TextAlignCenter, renderer.TextWrapPretty, &borderGray)
	rootStyled.AddChild(centerPrettyNode)

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
	footerStyled.Content = "Demonstrating TextWrap (None/Normal/Balanced/Pretty) and TextAlign (Left/Center/Right)"
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
