package renderer

import (
	"testing"

	"github.com/SCKelemen/color"
)

func TestNewStyle(t *testing.T) {
	s := NewStyle()

	if s == nil {
		t.Fatal("NewStyle() returned nil")
	}

	// Should be zero-valued
	if s.Bold {
		t.Error("Expected Bold to be false")
	}
	if s.Italic {
		t.Error("Expected Italic to be false")
	}
	if s.Foreground != nil {
		t.Error("Expected Foreground to be nil")
	}
}

func TestStyleWithForeground(t *testing.T) {
	c, _ := color.ParseColor("#FF0000")
	s := NewStyle().WithForeground(&c)

	if s.Foreground == nil {
		t.Fatal("Expected Foreground to be set")
	}
}

func TestStyleWithBackground(t *testing.T) {
	c, _ := color.ParseColor("#00FF00")
	s := NewStyle().WithBackground(&c)

	if s.Background == nil {
		t.Fatal("Expected Background to be set")
	}
}

func TestStyleWithBold(t *testing.T) {
	s := NewStyle().WithBold(true)

	if !s.Bold {
		t.Error("Expected Bold to be true")
	}
}

func TestStyleWithItalic(t *testing.T) {
	s := NewStyle().WithItalic(true)

	if !s.Italic {
		t.Error("Expected Italic to be true")
	}
}

func TestStyleWithUnderline(t *testing.T) {
	s := NewStyle().WithUnderline(true)

	if !s.Underline {
		t.Error("Expected Underline to be true")
	}
}

func TestStyleWithBorder(t *testing.T) {
	s := NewStyle().WithBorder(RoundedBorder)

	if s.Border == nil {
		t.Fatal("Expected Border to be set")
	}
	if !s.Border.Top || !s.Border.Right || !s.Border.Bottom || !s.Border.Left {
		t.Error("Expected all border sides to be true")
	}
	if s.Border.Chars != RoundedBorder {
		t.Error("Expected RoundedBorder chars")
	}
}

func TestStyleWithBorderColor(t *testing.T) {
	c, _ := color.ParseColor("#0000FF")
	s := NewStyle().WithBorderColor(&c)

	if s.BorderColor == nil {
		t.Fatal("Expected BorderColor to be set")
	}
}

func TestStyleWithTextOverflow(t *testing.T) {
	s := NewStyle().WithTextOverflow(TextOverflowEllipsis)

	if s.TextOverflow != TextOverflowEllipsis {
		t.Errorf("Expected TextOverflowEllipsis, got %d", s.TextOverflow)
	}
}

func TestStyleChaining(t *testing.T) {
	red, _ := color.ParseColor("#FF0000")
	blue, _ := color.ParseColor("#0000FF")

	s := NewStyle().
		WithForeground(&red).
		WithBackground(&blue).
		WithBold(true).
		WithItalic(true)

	if s.Foreground == nil {
		t.Error("Expected Foreground to be set")
	}
	if s.Background == nil {
		t.Error("Expected Background to be set")
	}
	if !s.Bold {
		t.Error("Expected Bold to be true")
	}
	if !s.Italic {
		t.Error("Expected Italic to be true")
	}
}

func TestBorderChars(t *testing.T) {
	tests := []struct {
		name   string
		border BorderChars
	}{
		{"RoundedBorder", RoundedBorder},
		{"DoubleBorder", DoubleBorder},
		{"ThickBorder", ThickBorder},
		{"NormalBorder", NormalBorder},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Check that all border characters are set
			if tt.border.TopLeft == 0 {
				t.Error("TopLeft should not be zero")
			}
			if tt.border.TopRight == 0 {
				t.Error("TopRight should not be zero")
			}
			if tt.border.BottomLeft == 0 {
				t.Error("BottomLeft should not be zero")
			}
			if tt.border.BottomRight == 0 {
				t.Error("BottomRight should not be zero")
			}
			if tt.border.Horizontal == 0 {
				t.Error("Horizontal should not be zero")
			}
			if tt.border.Vertical == 0 {
				t.Error("Vertical should not be zero")
			}
		})
	}
}

func TestTextOverflowConstants(t *testing.T) {
	// Just verify the constants exist and have distinct values
	modes := []TextOverflow{
		TextOverflowClip,
		TextOverflowEllipsis,
		TextOverflowEllipsisStart,
		TextOverflowEllipsisMiddle,
	}

	seen := make(map[TextOverflow]bool)
	for _, mode := range modes {
		if seen[mode] {
			t.Errorf("Duplicate TextOverflow value: %d", mode)
		}
		seen[mode] = true
	}
}

func TestTextWrapConstants(t *testing.T) {
	modes := []TextWrap{
		TextWrapNone,
		TextWrapNormal,
		TextWrapBalanced,
		TextWrapPretty,
	}

	seen := make(map[TextWrap]bool)
	for _, mode := range modes {
		if seen[mode] {
			t.Errorf("Duplicate TextWrap value: %d", mode)
		}
		seen[mode] = true
	}
}

func TestTextAlignConstants(t *testing.T) {
	aligns := []TextAlign{
		TextAlignLeft,
		TextAlignCenter,
		TextAlignRight,
		TextAlignJustify,
	}

	seen := make(map[TextAlign]bool)
	for _, align := range aligns {
		if seen[align] {
			t.Errorf("Duplicate TextAlign value: %d", align)
		}
		seen[align] = true
	}
}
