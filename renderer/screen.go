package renderer

import (
	"strings"

	"github.com/SCKelemen/layout"
	"github.com/SCKelemen/text"
)

// textMeasurer is the Unicode-aware text measurement library
var textMeasurer = text.NewTerminal()

// Cell represents a single character cell in the terminal
type Cell struct {
	// Content can be a single rune or a complete grapheme cluster (emoji sequence, etc.)
	Content string
	Style   *Style
}

// Screen represents the terminal screen buffer
type Screen struct {
	Width    int
	Height   int
	Cells    [][]Cell
	Previous [][]Cell
	renderer *ANSIRenderer
}

// NewScreen creates a new screen buffer
func NewScreen(width, height int) *Screen {
	return &Screen{
		Width:    width,
		Height:   height,
		Cells:    makeBuffer(width, height),
		Previous: makeBuffer(width, height),
		renderer: NewANSIRenderer(),
	}
}

// makeBuffer creates a 2D array of cells
func makeBuffer(width, height int) [][]Cell {
	buffer := make([][]Cell, height)
	for y := 0; y < height; y++ {
		buffer[y] = make([]Cell, width)
		for x := 0; x < width; x++ {
			buffer[y][x] = Cell{Content: " ", Style: nil}
		}
	}
	return buffer
}

// Resize changes the screen dimensions
func (s *Screen) Resize(width, height int) {
	if s.Width == width && s.Height == height {
		return
	}

	s.Width = width
	s.Height = height
	s.Cells = makeBuffer(width, height)
	s.Previous = makeBuffer(width, height)
}

// SetColorMode sets the color mode for rendering
func (s *Screen) SetColorMode(mode ColorMode) {
	s.renderer = NewANSIRendererWithMode(mode)
}

// Clear resets all cells to empty
func (s *Screen) Clear() {
	for y := 0; y < s.Height; y++ {
		for x := 0; x < s.Width; x++ {
			s.Cells[y][x] = Cell{Content: " ", Style: nil}
		}
	}
}

// SetCell sets a single cell with content (can be a rune or grapheme cluster)
func (s *Screen) SetCell(x, y int, content string, style *Style) {
	if x < 0 || x >= s.Width || y < 0 || y >= s.Height {
		return
	}
	s.Cells[y][x] = Cell{Content: content, Style: style}
}

// Render renders a styled node to the screen buffer
func (s *Screen) Render(node *StyledNode) {
	s.Clear()
	s.renderNodeWithOffset(node, 0, 0)
}

// renderNode recursively renders a node and its children (legacy method)
func (s *Screen) renderNode(node *StyledNode) {
	s.renderNodeWithOffset(node, 0, 0)
}

// renderNodeWithOffset recursively renders a node and its children with accumulated offsets
func (s *Screen) renderNodeWithOffset(node *StyledNode, offsetX, offsetY int) {
	if node == nil || node.Node == nil {
		return
	}

	// Calculate absolute position by adding parent offsets
	x := int(node.Node.Rect.X) + offsetX
	y := int(node.Node.Rect.Y) + offsetY
	w := int(node.Node.Rect.Width)
	h := int(node.Node.Rect.Height)

	// Render background if present
	if node.Style != nil && node.Style.Background != nil {
		s.renderBackground(x, y, w, h, node.Style)
	}

	// Render border if present
	if node.Style != nil && node.Style.Border != nil {
		s.renderBorder(x, y, w, h, node.Style)
	}

	// Render content
	if node.Content != "" {
		// Account for border offset
		contentX := x
		contentY := y
		contentW := w
		contentH := h

		if node.Style != nil && node.Style.Border != nil {
			contentX++
			contentY++
			contentW -= 2
			contentH -= 2
		}

		s.renderText(contentX, contentY, contentW, contentH, node.Content, node.Style)
	}

	// Render children with accumulated offsets
	for _, child := range node.Children {
		s.renderNodeWithOffset(child, x, y)
	}
}

// renderBorder renders a border around the specified rectangle
func (s *Screen) renderBorder(x, y, w, h int, style *Style) {
	if style.Border == nil {
		return
	}

	border := style.Border
	chars := border.Chars

	// Use border color if specified, otherwise use foreground color
	borderStyle := &Style{
		Foreground: style.BorderColor,
	}
	if borderStyle.Foreground == nil {
		borderStyle.Foreground = style.Foreground
	}

	// Top border
	if border.Top && y >= 0 && y < s.Height {
		if border.Left && x >= 0 && x < s.Width {
			s.SetCell(x, y, string(chars.TopLeft), borderStyle)
		}
		for i := 1; i < w-1; i++ {
			if x+i >= 0 && x+i < s.Width {
				s.SetCell(x+i, y, string(chars.Horizontal), borderStyle)
			}
		}
		if border.Right && x+w-1 >= 0 && x+w-1 < s.Width {
			s.SetCell(x+w-1, y, string(chars.TopRight), borderStyle)
		}
	}

	// Bottom border
	if border.Bottom && y+h-1 >= 0 && y+h-1 < s.Height {
		if border.Left && x >= 0 && x < s.Width {
			s.SetCell(x, y+h-1, string(chars.BottomLeft), borderStyle)
		}
		for i := 1; i < w-1; i++ {
			if x+i >= 0 && x+i < s.Width {
				s.SetCell(x+i, y+h-1, string(chars.Horizontal), borderStyle)
			}
		}
		if border.Right && x+w-1 >= 0 && x+w-1 < s.Width {
			s.SetCell(x+w-1, y+h-1, string(chars.BottomRight), borderStyle)
		}
	}

	// Left and right borders
	for i := 1; i < h-1; i++ {
		if border.Left && y+i >= 0 && y+i < s.Height && x >= 0 && x < s.Width {
			s.SetCell(x, y+i, string(chars.Vertical), borderStyle)
		}
		if border.Right && y+i >= 0 && y+i < s.Height && x+w-1 >= 0 && x+w-1 < s.Width {
			s.SetCell(x+w-1, y+i, string(chars.Vertical), borderStyle)
		}
	}
}

// renderBackground fills the rectangle with the background color
func (s *Screen) renderBackground(x, y, w, h int, style *Style) {
	if style == nil || style.Background == nil {
		return
	}

	// Create a background-only style (no foreground, no text attributes)
	bgStyle := &Style{
		Background: style.Background,
	}

	// Fill the entire rectangle with spaces
	for row := y; row < y+h; row++ {
		if row < 0 || row >= s.Height {
			continue
		}
		for col := x; col < x+w; col++ {
			if col < 0 || col >= s.Width {
				continue
			}
			s.SetCell(col, row, " ", bgStyle)
		}
	}
}

// renderText renders text within the specified rectangle
func (s *Screen) renderText(x, y, w, h int, content string, style *Style) {
	if content == "" {
		return
	}

	// Get lines based on wrapping mode
	var lines []text.Line
	wrapMode := TextWrapNone
	if style != nil {
		wrapMode = style.TextWrap
	}

	switch wrapMode {
	case TextWrapNone:
		// Manual line breaks only
		rawLines := strings.Split(content, "\n")
		lines = make([]text.Line, len(rawLines))
		for i, line := range rawLines {
			lines[i] = text.Line{
				Content: line,
				Width:   textMeasurer.Width(line),
			}
		}

	case TextWrapNormal, TextWrapBalanced, TextWrapPretty:
		// For wrapping modes, we need to handle manual line breaks (paragraphs) first
		// Split on \n to get paragraphs, wrap each separately, then combine
		paragraphs := strings.Split(content, "\n")
		lines = make([]text.Line, 0)

		for _, para := range paragraphs {
			if para == "" {
				// Empty line - preserve as blank line
				lines = append(lines, text.Line{
					Content: "",
					Width:   0,
				})
				continue
			}

			var paraLines []text.Line
			switch wrapMode {
			case TextWrapNormal:
				paraLines = textMeasurer.Wrap(para, text.WrapOptions{
					MaxWidth:   float64(w),
					BreakWords: false,
				})
			case TextWrapBalanced:
				paraLines = textMeasurer.WrapBalanced(para, float64(w))
			case TextWrapPretty:
				paraLines = textMeasurer.WrapKnuthPlass(para, text.KnuthPlassOptions{
					MaxWidth:  float64(w),
					Tolerance: 1.0,
				})
			}
			lines = append(lines, paraLines...)
		}
	}

	// Render each line
	for lineIdx, line := range lines {
		if lineIdx >= h {
			break
		}

		row := y + lineIdx
		if row < 0 || row >= s.Height {
			continue
		}

		lineText := line.Content
		lineWidth := line.Width

		// Apply ellipsis if line overflows
		if style != nil && lineWidth > float64(w) {
			switch style.TextOverflow {
			case TextOverflowEllipsis:
				lineText = textMeasurer.ElideEndWith(lineText, float64(w), "…")
			case TextOverflowEllipsisStart:
				lineText = textMeasurer.ElideStartWith(lineText, float64(w), "…")
			case TextOverflowEllipsisMiddle:
				lineText = textMeasurer.ElideWith(lineText, float64(w), "…")
			}
			lineWidth = textMeasurer.Width(lineText)
		}

		// Apply justify alignment (if not last line and text is short enough)
		isLastLine := lineIdx == len(lines)-1
		if style != nil && style.TextAlign == TextAlignJustify && !isLastLine && lineWidth < float64(w) {
			// Only justify if line has multiple words
			if strings.Contains(lineText, " ") {
				lineText = textMeasurer.JustifyText(lineText, float64(w), text.TextJustifyInterWord)
				lineWidth = textMeasurer.Width(lineText)
			}
		}

		// Calculate starting column based on alignment
		col := x
		if style != nil {
			switch style.TextAlign {
			case TextAlignCenter:
				col = x + (w-int(lineWidth))/2
			case TextAlignRight:
				col = x + w - int(lineWidth)
			case TextAlignJustify:
				// Justify was already applied above, so left-align the result
				col = x
			case TextAlignLeft:
				col = x
			}
		}

		// Render the line with proper grapheme cluster handling
		graphemes := textMeasurer.Graphemes(lineText)

		for _, grapheme := range graphemes {
			// Measure grapheme width
			graphemeWidth := int(textMeasurer.Width(grapheme))

			// Check if grapheme fits in the remaining space
			if col+graphemeWidth > x+w || col >= s.Width {
				break
			}

			// Render the grapheme cluster
			// For single-rune graphemes (most characters), this outputs one character
			// For multi-rune graphemes (emoji sequences, combining marks), we store
			// the entire sequence in the cell and it will be output correctly
			// Note: Complex emoji sequences may not render correctly in all terminals
			if len(grapheme) > 0 && col >= 0 {
				// Store the complete grapheme cluster (handles emoji sequences correctly)
				s.SetCell(col, row, grapheme, style)

				// For wide characters/graphemes (width=2), mark the second column
				// This prevents other content from overlapping
				if graphemeWidth == 2 && col+1 < s.Width {
					s.SetCell(col+1, row, " ", style)
				}
			}

			col += graphemeWidth
		}
	}
}

// String converts the screen buffer to a string with ANSI codes
func (s *Screen) String() string {
	var buf strings.Builder

	// Start with cursor at top-left
	buf.WriteString(s.renderer.MoveCursor(0, 0))

	var lastStyle *Style
	for y := 0; y < s.Height; y++ {
		for x := 0; x < s.Width; x++ {
			cell := s.Cells[y][x]

			// Only output ANSI codes when style changes
			if !stylesEqual(cell.Style, lastStyle) {
				buf.WriteString(s.renderer.Reset())
				if cell.Style != nil {
					buf.WriteString(s.renderer.RenderStyle(cell.Style))
				}
				lastStyle = cell.Style
			}

			// Output the cell content (can be a single character or emoji sequence)
			buf.WriteString(cell.Content)
		}
		if y < s.Height-1 {
			buf.WriteString("\n")
		}
	}

	buf.WriteString(s.renderer.Reset())
	return buf.String()
}

// stylesEqual checks if two styles are equal
func stylesEqual(a, b *Style) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// For now, use pointer equality
	// We could implement a more sophisticated comparison if needed
	return a == b
}

// StyledNode wraps a layout.Node with visual styling
type StyledNode struct {
	*layout.Node
	Style    *Style
	Content  string
	Children []*StyledNode
}

// NewStyledNode creates a new styled node
func NewStyledNode(node *layout.Node, style *Style) *StyledNode {
	return &StyledNode{
		Node:  node,
		Style: style,
	}
}

// AddChild adds a child node
func (n *StyledNode) AddChild(child *StyledNode) {
	n.Children = append(n.Children, child)
	if n.Node != nil && child.Node != nil {
		n.Node.Children = append(n.Node.Children, child.Node)
	}
}
