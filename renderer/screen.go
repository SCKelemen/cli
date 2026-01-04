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
	s.renderNode(node)
}

// renderNode recursively renders a node and its children
func (s *Screen) renderNode(node *StyledNode) {
	if node == nil || node.Node == nil {
		return
	}

	x := int(node.Node.Rect.X)
	y := int(node.Node.Rect.Y)
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

	// Render children
	for _, child := range node.Children {
		s.renderNode(child)
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

	lines := strings.Split(content, "\n")
	for lineIdx, line := range lines {
		if lineIdx >= h {
			break
		}

		row := y + lineIdx
		if row < 0 || row >= s.Height {
			continue
		}

		// Measure line width with proper Unicode handling
		lineWidth := textMeasurer.Width(line)

		// Apply ellipsis if line overflows
		if style != nil && lineWidth > float64(w) {
			switch style.TextOverflow {
			case TextOverflowEllipsis:
				line = textMeasurer.ElideEndWith(line, float64(w), "…")
			case TextOverflowEllipsisStart:
				line = textMeasurer.ElideStartWith(line, float64(w), "…")
			case TextOverflowEllipsisMiddle:
				line = textMeasurer.ElideWith(line, float64(w), "…")
			}
		}

		// Render the line with proper grapheme cluster handling
		col := x
		graphemes := textMeasurer.Graphemes(line)

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
