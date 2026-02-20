package main

import (
	"fmt"

	"github.com/SCKelemen/color"
	"github.com/SCKelemen/layout"
	"github.com/SCKelemen/cli/components"
	"github.com/SCKelemen/cli/renderer"
)

func main() {
	width, height := 70, 25

	screen := renderer.NewScreen(width, height)
	root := buildEllipsisDemo(width, height)

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

func buildEllipsisDemo(width, height int) *renderer.StyledNode {
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
	bgBlue, _ := color.ParseColor("oklch(0.5 0.2 240)")
	borderBlue, _ := color.ParseColor("oklch(0.7 0.2 240)")

	header := components.NewMessageBlock("Text Ellipsis Modes Demo")
	headerStyle := &renderer.Style{
		Foreground:  &fgWhite,
		Background:  &bgBlue,
		Bold:        true,
		BorderColor: &borderBlue,
	}
	headerStyle.WithBorder(renderer.RoundedBorder)
	headerNode := header.ToStyledNode()
	headerNode.Style = headerStyle
	rootStyled.AddChild(headerNode)

	borderGray, _ := color.ParseColor("#5A5A5A")
	longText := "This is a very long line of text that should be truncated with an ellipsis when it overflows the available width in the terminal."

	// End ellipsis (default)
	endMsg := components.NewMessageBlock("End: " + longText)
	endMsg = endMsg.WithBorderColor(&borderGray)
	endNode := endMsg.ToStyledNode()
	endNode.Style.TextOverflow = renderer.TextOverflowEllipsis
	rootStyled.AddChild(endNode)

	// Start ellipsis
	startMsg := components.NewMessageBlock("Start: " + longText)
	startMsg = startMsg.WithBorderColor(&borderGray)
	startNode := startMsg.ToStyledNode()
	startNode.Style.TextOverflow = renderer.TextOverflowEllipsisStart
	rootStyled.AddChild(startNode)

	// Middle ellipsis
	middleMsg := components.NewMessageBlock("Middle: " + longText)
	middleMsg = middleMsg.WithBorderColor(&borderGray)
	middleNode := middleMsg.ToStyledNode()
	middleNode.Style.TextOverflow = renderer.TextOverflowEllipsisMiddle
	rootStyled.AddChild(middleNode)

	// Path example with middle ellipsis
	pathText := "/usr/local/share/applications/really/deeply/nested/directory/structure/myapp.desktop"
	pathMsg := components.NewMessageBlock("Path: " + pathText)
	pathMsg = pathMsg.WithBorderColor(&borderGray)
	pathNode := pathMsg.ToStyledNode()
	pathNode.Style.TextOverflow = renderer.TextOverflowEllipsisMiddle
	rootStyled.AddChild(pathNode)

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
	}
	footerStyled := renderer.NewStyledNode(footerNode, footerStyle)
	footerStyled.Content = "Demonstrating TextOverflowEllipsis, TextOverflowEllipsisStart, and TextOverflowEllipsisMiddle"
	rootStyled.AddChild(footerStyled)

	return rootStyled
}
