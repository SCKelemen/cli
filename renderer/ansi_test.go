package renderer

import (
	"strings"
	"testing"

	"github.com/SCKelemen/color"
)

func TestNewANSIRenderer(t *testing.T) {
	r := NewANSIRenderer()

	if r == nil {
		t.Fatal("NewANSIRenderer() returned nil")
	}
}

func TestNewANSIRendererWithMode(t *testing.T) {
	modes := []ColorMode{
		ColorModeNone,
		ColorMode16,
		ColorMode256,
		ColorModeTrueColor,
	}

	for _, mode := range modes {
		r := NewANSIRendererWithMode(mode)
		if r == nil {
			t.Errorf("NewANSIRendererWithMode(%d) returned nil", mode)
		}
	}
}

func TestRenderStyleBold(t *testing.T) {
	r := NewANSIRenderer()
	style := &Style{Bold: true}

	output := r.RenderStyle(style)

	if !strings.Contains(output, "\x1b[1m") {
		t.Error("Expected bold ANSI code \\x1b[1m")
	}
}

func TestRenderStyleItalic(t *testing.T) {
	r := NewANSIRenderer()
	style := &Style{Italic: true}

	output := r.RenderStyle(style)

	if !strings.Contains(output, "\x1b[3m") {
		t.Error("Expected italic ANSI code \\x1b[3m")
	}
}

func TestRenderStyleUnderline(t *testing.T) {
	r := NewANSIRenderer()
	style := &Style{Underline: true}

	output := r.RenderStyle(style)

	if !strings.Contains(output, "\x1b[4m") {
		t.Error("Expected underline ANSI code \\x1b[4m")
	}
}

func TestRenderStyleDim(t *testing.T) {
	r := NewANSIRenderer()
	style := &Style{Dim: true}

	output := r.RenderStyle(style)

	if !strings.Contains(output, "\x1b[2m") {
		t.Error("Expected dim ANSI code \\x1b[2m")
	}
}

func TestRenderStyleBlink(t *testing.T) {
	r := NewANSIRenderer()
	style := &Style{Blink: true}

	output := r.RenderStyle(style)

	if !strings.Contains(output, "\x1b[5m") {
		t.Error("Expected blink ANSI code \\x1b[5m")
	}
}

func TestRenderStyleReverse(t *testing.T) {
	r := NewANSIRenderer()
	style := &Style{Reverse: true}

	output := r.RenderStyle(style)

	if !strings.Contains(output, "\x1b[7m") {
		t.Error("Expected reverse ANSI code \\x1b[7m")
	}
}

func TestRenderStyleStrikethrough(t *testing.T) {
	r := NewANSIRenderer()
	style := &Style{Strikethrough: true}

	output := r.RenderStyle(style)

	if !strings.Contains(output, "\x1b[9m") {
		t.Error("Expected strikethrough ANSI code \\x1b[9m")
	}
}

func TestRenderStyleWithForeground(t *testing.T) {
	r := NewANSIRendererWithMode(ColorModeTrueColor)
	red, _ := color.ParseColor("#FF0000")
	style := &Style{Foreground: &red}

	output := r.RenderStyle(style)

	// Should contain some foreground color code
	if !strings.Contains(output, "\x1b[") {
		t.Error("Expected ANSI escape sequence")
	}
}

func TestRenderStyleWithBackground(t *testing.T) {
	r := NewANSIRendererWithMode(ColorModeTrueColor)
	blue, _ := color.ParseColor("#0000FF")
	style := &Style{Background: &blue}

	output := r.RenderStyle(style)

	// Should contain some background color code
	if !strings.Contains(output, "\x1b[") {
		t.Error("Expected ANSI escape sequence")
	}
}

func TestRenderStyleMultipleAttributes(t *testing.T) {
	r := NewANSIRendererWithMode(ColorModeTrueColor)
	red, _ := color.ParseColor("#FF0000")
	style := &Style{
		Foreground: &red,
		Bold:       true,
		Underline:  true,
	}

	output := r.RenderStyle(style)

	if !strings.Contains(output, "1") {
		t.Error("Expected bold code")
	}
	if !strings.Contains(output, "4") {
		t.Error("Expected underline code")
	}
}

func TestRenderStyleNil(t *testing.T) {
	r := NewANSIRenderer()

	output := r.RenderStyle(nil)

	// Should return empty string or just reset
	if len(output) > 10 {
		t.Error("Expected minimal output for nil style")
	}
}

func TestReset(t *testing.T) {
	r := NewANSIRenderer()

	output := r.Reset()

	if !strings.Contains(output, "\x1b[0m") {
		t.Error("Expected reset ANSI code \\x1b[0m")
	}
}

func TestMoveCursor(t *testing.T) {
	r := NewANSIRenderer()

	tests := []struct {
		x, y int
	}{
		{0, 0},
		{10, 5},
		{80, 24},
	}

	for _, tt := range tests {
		output := r.MoveCursor(tt.x, tt.y)
		if !strings.Contains(output, "\x1b[") {
			t.Errorf("MoveCursor(%d, %d) should contain ANSI escape sequence", tt.x, tt.y)
		}
	}
}

func TestClearScreen(t *testing.T) {
	r := NewANSIRenderer()

	output := r.ClearScreen()

	if !strings.Contains(output, "\x1b[") {
		t.Error("Expected ANSI escape sequence for clear screen")
	}
}

func TestHideCursor(t *testing.T) {
	r := NewANSIRenderer()

	output := r.HideCursor()

	if !strings.Contains(output, "\x1b[?25l") {
		t.Error("Expected hide cursor ANSI code \\x1b[?25l")
	}
}

func TestShowCursor(t *testing.T) {
	r := NewANSIRenderer()

	output := r.ShowCursor()

	if !strings.Contains(output, "\x1b[?25h") {
		t.Error("Expected show cursor ANSI code \\x1b[?25h")
	}
}

func TestColorModeConstants(t *testing.T) {
	modes := []ColorMode{
		ColorModeNone,
		ColorMode16,
		ColorMode256,
		ColorModeTrueColor,
	}

	seen := make(map[ColorMode]bool)
	for _, mode := range modes {
		if seen[mode] {
			t.Errorf("Duplicate ColorMode value: %d", mode)
		}
		seen[mode] = true
	}
}
