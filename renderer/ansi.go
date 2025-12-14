package renderer

import (
	"fmt"
	"strings"
)

// ANSIRenderer converts styles to ANSI escape codes
type ANSIRenderer struct{}

// NewANSIRenderer creates a new ANSI renderer
func NewANSIRenderer() *ANSIRenderer {
	return &ANSIRenderer{}
}

// Reset returns the ANSI reset sequence
func (r *ANSIRenderer) Reset() string {
	return "\x1b[0m"
}

// RenderStyle converts a style to ANSI escape codes
func (r *ANSIRenderer) RenderStyle(s *Style) string {
	if s == nil {
		return ""
	}

	var codes []string

	// Text attributes
	if s.Bold {
		codes = append(codes, "1")
	}
	if s.Dim {
		codes = append(codes, "2")
	}
	if s.Italic {
		codes = append(codes, "3")
	}
	if s.Underline {
		codes = append(codes, "4")
	}
	if s.Blink {
		codes = append(codes, "5")
	}
	if s.Reverse {
		codes = append(codes, "7")
	}
	if s.Strikethrough {
		codes = append(codes, "9")
	}

	// Foreground color
	if s.Foreground != nil {
		r, g, b, _ := (*s.Foreground).RGBA()
		codes = append(codes, fmt.Sprintf("38;2;%d;%d;%d",
			int(r*255), int(g*255), int(b*255)))
	}

	// Background color
	if s.Background != nil {
		r, g, b, _ := (*s.Background).RGBA()
		codes = append(codes, fmt.Sprintf("48;2;%d;%d;%d",
			int(r*255), int(g*255), int(b*255)))
	}

	if len(codes) == 0 {
		return ""
	}

	return "\x1b[" + strings.Join(codes, ";") + "m"
}

// MoveCursor moves the cursor to the specified position (1-indexed)
func (r *ANSIRenderer) MoveCursor(x, y int) string {
	return fmt.Sprintf("\x1b[%d;%dH", y+1, x+1)
}

// ClearScreen clears the entire screen
func (r *ANSIRenderer) ClearScreen() string {
	return "\x1b[2J\x1b[H"
}

// HideCursor hides the cursor
func (r *ANSIRenderer) HideCursor() string {
	return "\x1b[?25l"
}

// ShowCursor shows the cursor
func (r *ANSIRenderer) ShowCursor() string {
	return "\x1b[?25h"
}

// EnterAltScreen enters the alternate screen buffer
func (r *ANSIRenderer) EnterAltScreen() string {
	return "\x1b[?1049h"
}

// ExitAltScreen exits the alternate screen buffer
func (r *ANSIRenderer) ExitAltScreen() string {
	return "\x1b[?1049l"
}
