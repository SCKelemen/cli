package renderer

import (
	"strings"

	"github.com/SCKelemen/layout"
	"github.com/mattn/go-runewidth"
)

// Cell represents a single character cell in the terminal
type Cell struct {
	Char  rune
	Style *Style
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
			buffer[y][x] = Cell{Char: ' ', Style: nil}
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
			s.Cells[y][x] = Cell{Char: ' ', Style: nil}
		}
	}
}

// SetCell sets a single cell
func (s *Screen) SetCell(x, y int, char rune, style *Style) {
	if x < 0 || x >= s.Width || y < 0 || y >= s.Height {
		return
	}
	s.Cells[y][x] = Cell{Char: char, Style: style}
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
			s.SetCell(x, y, chars.TopLeft, borderStyle)
		}
		for i := 1; i < w-1; i++ {
			if x+i >= 0 && x+i < s.Width {
				s.SetCell(x+i, y, chars.Horizontal, borderStyle)
			}
		}
		if border.Right && x+w-1 >= 0 && x+w-1 < s.Width {
			s.SetCell(x+w-1, y, chars.TopRight, borderStyle)
		}
	}

	// Bottom border
	if border.Bottom && y+h-1 >= 0 && y+h-1 < s.Height {
		if border.Left && x >= 0 && x < s.Width {
			s.SetCell(x, y+h-1, chars.BottomLeft, borderStyle)
		}
		for i := 1; i < w-1; i++ {
			if x+i >= 0 && x+i < s.Width {
				s.SetCell(x+i, y+h-1, chars.Horizontal, borderStyle)
			}
		}
		if border.Right && x+w-1 >= 0 && x+w-1 < s.Width {
			s.SetCell(x+w-1, y+h-1, chars.BottomRight, borderStyle)
		}
	}

	// Left and right borders
	for i := 1; i < h-1; i++ {
		if border.Left && y+i >= 0 && y+i < s.Height && x >= 0 && x < s.Width {
			s.SetCell(x, y+i, chars.Vertical, borderStyle)
		}
		if border.Right && y+i >= 0 && y+i < s.Height && x+w-1 >= 0 && x+w-1 < s.Width {
			s.SetCell(x+w-1, y+i, chars.Vertical, borderStyle)
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
			s.SetCell(col, row, ' ', bgStyle)
		}
	}
}

// renderText renders text within the specified rectangle
func (s *Screen) renderText(x, y, w, h int, text string, style *Style) {
	if text == "" {
		return
	}

	lines := strings.Split(text, "\n")
	for lineIdx, line := range lines {
		if lineIdx >= h {
			break
		}

		row := y + lineIdx
		if row < 0 || row >= s.Height {
			continue
		}

		col := x
		for _, char := range line {
			if col >= x+w || col >= s.Width {
				break
			}
			if col >= 0 {
				s.SetCell(col, row, char, style)
			}
			// Account for character display width (some chars take 2 columns)
			col += runewidth.RuneWidth(char)
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

			buf.WriteRune(cell.Char)
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
