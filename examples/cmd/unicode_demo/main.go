package main

import (
	"fmt"

	"github.com/SCKelemen/cli/components"
	"github.com/SCKelemen/cli/renderer"
	"github.com/SCKelemen/color"
	"github.com/SCKelemen/layout"
)

func main() {
	width, height := 70, 35

	screen := renderer.NewScreen(width, height)
	root := buildUnicodeTest(width, height)

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

func buildUnicodeTest(width, height int) *renderer.StyledNode {
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

	header := components.NewMessageBlock("Unicode Text Measurement Test")
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

	// Test cases
	testCases := []struct {
		label string
		text  string
	}{
		{"ASCII:", "Hello World!"},
		{"Emoji:", "Hello 👋 World 🌍"},
		{"Emoji Seq:", "Skin tones: 👋🏻 👋🏽 👋🏿"},
		{"Flags:", "Flags: 🇺🇸 🇯🇵 🇬🇧 🇫🇷"},
		{"CJK:", "日本語 中文 한글"},
		{"Mixed:", "Hello 世界 👋 🌍"},
		{"Combining:", "Café naïve"},
		{"ZWJ Emoji:", "Family: 👨‍👩‍👧‍👦"},
	}

	borderGray, _ := color.ParseColor("#5A5A5A")

	for _, tc := range testCases {
		msg := components.NewMessageBlock(fmt.Sprintf("%s %s", tc.label, tc.text))
		msg = msg.WithBorderColor(&borderGray)
		rootStyled.AddChild(msg.ToStyledNode())
	}

	// Ellipsis test - long text that should be truncated
	longText := "This is a very long line of text that should be truncated with an ellipsis when it overflows the available width in the terminal. The text library's ElideEnd function handles this properly."

	ellipsisMsg := components.NewMessageBlock(longText)
	ellipsisMsg = ellipsisMsg.WithBorderColor(&borderGray)
	ellipsisNode := ellipsisMsg.ToStyledNode()
	ellipsisNode.Style.TextOverflow = renderer.TextOverflowEllipsis
	rootStyled.AddChild(ellipsisNode)

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
	footerStyled.Content = "Using SCKelemen/text for proper Unicode measurement"
	rootStyled.AddChild(footerStyled)

	return rootStyled
}
