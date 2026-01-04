package renderer

import (
	"strings"
	"testing"

	"github.com/SCKelemen/color"
	"github.com/SCKelemen/layout"
)

func TestNewScreen(t *testing.T) {
	s := NewScreen(80, 24)

	if s.Width != 80 {
		t.Errorf("Expected width 80, got %d", s.Width)
	}
	if s.Height != 24 {
		t.Errorf("Expected height 24, got %d", s.Height)
	}
	if len(s.Cells) != 24 {
		t.Errorf("Expected 24 rows, got %d", len(s.Cells))
	}
	if len(s.Cells[0]) != 80 {
		t.Errorf("Expected 80 columns, got %d", len(s.Cells[0]))
	}

	// Check cells are initialized with spaces
	for y := 0; y < s.Height; y++ {
		for x := 0; x < s.Width; x++ {
			if s.Cells[y][x].Content != " " {
				t.Errorf("Expected cell [%d][%d] to be space, got %q", y, x, s.Cells[y][x].Content)
			}
		}
	}
}

func TestScreenResize(t *testing.T) {
	s := NewScreen(80, 24)
	s.Resize(100, 30)

	if s.Width != 100 {
		t.Errorf("Expected width 100, got %d", s.Width)
	}
	if s.Height != 30 {
		t.Errorf("Expected height 30, got %d", s.Height)
	}
	if len(s.Cells) != 30 {
		t.Errorf("Expected 30 rows, got %d", len(s.Cells))
	}
	if len(s.Cells[0]) != 100 {
		t.Errorf("Expected 100 columns, got %d", len(s.Cells[0]))
	}
}

func TestScreenResizeSameSize(t *testing.T) {
	s := NewScreen(80, 24)
	originalWidth := s.Width
	originalHeight := s.Height

	s.Resize(80, 24) // Same size

	// Should be no-op - dimensions unchanged
	if s.Width != originalWidth || s.Height != originalHeight {
		t.Error("Resize to same size should not change dimensions")
	}
}

func TestScreenClear(t *testing.T) {
	s := NewScreen(10, 5)

	// Set some cells
	s.SetCell(0, 0, "A", nil)
	s.SetCell(5, 2, "B", nil)

	s.Clear()

	// All cells should be spaces again
	for y := 0; y < s.Height; y++ {
		for x := 0; x < s.Width; x++ {
			if s.Cells[y][x].Content != " " {
				t.Errorf("Expected cell [%d][%d] to be space after clear, got %q", y, x, s.Cells[y][x].Content)
			}
		}
	}
}

func TestSetCell(t *testing.T) {
	s := NewScreen(10, 5)
	style := &Style{Bold: true}

	s.SetCell(5, 2, "X", style)

	if s.Cells[2][5].Content != "X" {
		t.Errorf("Expected cell content 'X', got %q", s.Cells[2][5].Content)
	}
	if s.Cells[2][5].Style != style {
		t.Error("Expected cell style to match")
	}
}

func TestSetCellOutOfBounds(t *testing.T) {
	s := NewScreen(10, 5)

	// These should not panic
	s.SetCell(-1, 0, "A", nil)
	s.SetCell(0, -1, "A", nil)
	s.SetCell(10, 0, "A", nil)
	s.SetCell(0, 5, "A", nil)

	// Verify cells are unchanged
	for y := 0; y < s.Height; y++ {
		for x := 0; x < s.Width; x++ {
			if s.Cells[y][x].Content != " " {
				t.Errorf("Out of bounds SetCell modified cell [%d][%d]", y, x)
			}
		}
	}
}

func TestRenderSimpleText(t *testing.T) {
	s := NewScreen(20, 5)

	node := &layout.Node{
		Rect: layout.Rect{X: 0, Y: 0, Width: 20, Height: 5},
	}
	styledNode := NewStyledNode(node, nil)
	styledNode.Content = "Hello"

	s.Render(styledNode)

	// Check that "Hello" is rendered at the top-left
	expected := "Hello"
	for i, ch := range expected {
		if s.Cells[0][i].Content != string(ch) {
			t.Errorf("Expected cell [0][%d] to be %q, got %q", i, string(ch), s.Cells[0][i].Content)
		}
	}
}

func TestRenderWithNewlines(t *testing.T) {
	s := NewScreen(20, 5)

	node := &layout.Node{
		Rect: layout.Rect{X: 0, Y: 0, Width: 20, Height: 5},
	}
	styledNode := NewStyledNode(node, nil)
	styledNode.Content = "Line 1\nLine 2"

	s.Render(styledNode)

	// Check first line
	line1 := "Line 1"
	for i, ch := range line1 {
		if s.Cells[0][i].Content != string(ch) {
			t.Errorf("Expected cell [0][%d] to be %q, got %q", i, string(ch), s.Cells[0][i].Content)
		}
	}

	// Check second line
	line2 := "Line 2"
	for i, ch := range line2 {
		if s.Cells[1][i].Content != string(ch) {
			t.Errorf("Expected cell [1][%d] to be %q, got %q", i, string(ch), s.Cells[1][i].Content)
		}
	}
}

func TestRenderWithEmoji(t *testing.T) {
	s := NewScreen(20, 5)

	node := &layout.Node{
		Rect: layout.Rect{X: 0, Y: 0, Width: 20, Height: 5},
	}
	styledNode := NewStyledNode(node, nil)
	styledNode.Content = "Hi 👋"

	s.Render(styledNode)

	// Check that emoji is stored as complete grapheme cluster
	if !strings.Contains(s.Cells[0][3].Content, "👋") {
		t.Errorf("Expected cell to contain emoji, got %q", s.Cells[0][3].Content)
	}

	// Wide character should occupy next cell with space
	if s.Cells[0][4].Content != " " {
		t.Errorf("Expected cell after wide emoji to be space, got %q", s.Cells[0][4].Content)
	}
}

func TestTextAlignLeft(t *testing.T) {
	s := NewScreen(20, 3)

	node := &layout.Node{
		Rect: layout.Rect{X: 0, Y: 0, Width: 20, Height: 3},
	}
	style := &Style{TextAlign: TextAlignLeft}
	styledNode := NewStyledNode(node, style)
	styledNode.Content = "Hello"

	s.Render(styledNode)

	// Should start at column 0
	if s.Cells[0][0].Content != "H" {
		t.Errorf("Left-aligned text should start at column 0")
	}
}

func TestTextAlignCenter(t *testing.T) {
	s := NewScreen(20, 3)

	node := &layout.Node{
		Rect: layout.Rect{X: 0, Y: 0, Width: 20, Height: 3},
	}
	style := &Style{TextAlign: TextAlignCenter}
	styledNode := NewStyledNode(node, style)
	styledNode.Content = "Hi" // 2 chars, should be centered in 20-char width

	s.Render(styledNode)

	// Width=20, text=2, so (20-2)/2 = 9 spaces before
	if s.Cells[0][9].Content != "H" {
		t.Errorf("Center-aligned text should start at column 9, found 'H' at different position")
	}
}

func TestTextAlignRight(t *testing.T) {
	s := NewScreen(20, 3)

	node := &layout.Node{
		Rect: layout.Rect{X: 0, Y: 0, Width: 20, Height: 3},
	}
	style := &Style{TextAlign: TextAlignRight}
	styledNode := NewStyledNode(node, style)
	styledNode.Content = "Hi" // 2 chars

	s.Render(styledNode)

	// Should end at column 19, so start at 20-2=18
	if s.Cells[0][18].Content != "H" {
		t.Errorf("Right-aligned text should start at column 18")
	}
}

func TestTextWrapNormal(t *testing.T) {
	s := NewScreen(10, 5)

	node := &layout.Node{
		Rect: layout.Rect{X: 0, Y: 0, Width: 10, Height: 5},
	}
	style := &Style{TextWrap: TextWrapNormal}
	styledNode := NewStyledNode(node, style)
	styledNode.Content = "This is a long text that should wrap"

	s.Render(styledNode)

	// Text should wrap to multiple lines
	row0HasContent := false
	row1HasContent := false

	for x := 0; x < s.Width; x++ {
		if s.Cells[0][x].Content != " " {
			row0HasContent = true
		}
		if s.Cells[1][x].Content != " " {
			row1HasContent = true
		}
	}

	if !row0HasContent {
		t.Error("Expected row 0 to have content")
	}
	if !row1HasContent {
		t.Error("Expected row 1 to have content (text should wrap)")
	}
}

func TestTextOverflowEllipsis(t *testing.T) {
	s := NewScreen(10, 3)

	node := &layout.Node{
		Rect: layout.Rect{X: 0, Y: 0, Width: 10, Height: 3},
	}
	style := &Style{
		TextWrap:     TextWrapNone,
		TextOverflow: TextOverflowEllipsis,
	}
	styledNode := NewStyledNode(node, style)
	styledNode.Content = "This is a very long line"

	s.Render(styledNode)

	// Should have ellipsis at the end
	hasEllipsis := false
	for x := 0; x < s.Width; x++ {
		if strings.Contains(s.Cells[0][x].Content, "…") {
			hasEllipsis = true
			break
		}
	}

	if !hasEllipsis {
		t.Error("Expected ellipsis in overflowing text")
	}
}

func TestRenderWithBorder(t *testing.T) {
	s := NewScreen(10, 5)

	node := &layout.Node{
		Rect: layout.Rect{X: 0, Y: 0, Width: 10, Height: 5},
	}
	blue, _ := color.ParseColor("#0000FF")
	style := &Style{
		BorderColor: &blue,
	}
	style.WithBorder(NormalBorder)
	styledNode := NewStyledNode(node, style)
	styledNode.Content = "Test"

	s.Render(styledNode)

	// Check corners
	topLeft := s.Cells[0][0].Content
	topRight := s.Cells[0][9].Content
	bottomLeft := s.Cells[4][0].Content
	bottomRight := s.Cells[4][9].Content

	if topLeft != "┌" {
		t.Errorf("Expected top-left corner '┌', got %q", topLeft)
	}
	if topRight != "┐" {
		t.Errorf("Expected top-right corner '┐', got %q", topRight)
	}
	if bottomLeft != "└" {
		t.Errorf("Expected bottom-left corner '└', got %q", bottomLeft)
	}
	if bottomRight != "┘" {
		t.Errorf("Expected bottom-right corner '┘', got %q", bottomRight)
	}
}

func TestStringOutput(t *testing.T) {
	s := NewScreen(5, 2)
	s.SetCell(0, 0, "A", nil)
	s.SetCell(1, 0, "B", nil)
	s.SetCell(0, 1, "C", nil)

	output := s.String()

	// Should contain ANSI codes and characters
	if !strings.Contains(output, "A") {
		t.Error("Output should contain 'A'")
	}
	if !strings.Contains(output, "B") {
		t.Error("Output should contain 'B'")
	}
	if !strings.Contains(output, "C") {
		t.Error("Output should contain 'C'")
	}

	// Should contain newline between rows
	if !strings.Contains(output, "\n") {
		t.Error("Output should contain newline")
	}
}
