package renderer

import (
	"fmt"
	"strings"

	"github.com/SCKelemen/color"
)

// ANSIRenderer converts styles to ANSI escape codes
type ANSIRenderer struct {
	ColorMode ColorMode
}

// NewANSIRenderer creates a new ANSI renderer with detected capabilities
func NewANSIRenderer() *ANSIRenderer {
	caps := DetectCapabilities()
	return &ANSIRenderer{
		ColorMode: caps.ColorMode,
	}
}

// NewANSIRendererWithMode creates a new ANSI renderer with specified color mode
func NewANSIRendererWithMode(mode ColorMode) *ANSIRenderer {
	return &ANSIRenderer{
		ColorMode: mode,
	}
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
		colorCode := r.renderColor(s.Foreground, true)
		if colorCode != "" {
			codes = append(codes, colorCode)
		}
	}

	// Background color
	if s.Background != nil {
		colorCode := r.renderColor(s.Background, false)
		if colorCode != "" {
			codes = append(codes, colorCode)
		}
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

// renderColor renders a color based on terminal capabilities
func (r *ANSIRenderer) renderColor(c *color.Color, foreground bool) string {
	if c == nil {
		return ""
	}

	red, green, blue, _ := (*c).RGBA()
	r8 := int(red * 255)
	g8 := int(green * 255)
	b8 := int(blue * 255)

	prefix := "38" // Foreground
	if !foreground {
		prefix = "48" // Background
	}

	switch r.ColorMode {
	case ColorModeNone:
		return ""

	case ColorMode16:
		// Convert to closest 16-color ANSI code
		ansi16 := rgbToANSI16(r8, g8, b8)
		if foreground {
			return fmt.Sprintf("%d", ansi16)
		}
		return fmt.Sprintf("%d", ansi16+10) // Background = foreground + 10

	case ColorMode256:
		// Convert to 256-color palette
		ansi256 := rgbToANSI256(r8, g8, b8)
		return fmt.Sprintf("%s;5;%d", prefix, ansi256)

	case ColorModeTrueColor:
		// Use full 24-bit RGB
		return fmt.Sprintf("%s;2;%d;%d;%d", prefix, r8, g8, b8)

	default:
		return ""
	}
}

// rgbToANSI16 converts RGB to the closest 16-color ANSI code
func rgbToANSI16(r, g, b int) int {
	// Calculate brightness
	brightness := (r + g + b) / 3

	// Determine base color based on dominant channel
	if r > g && r > b {
		if brightness > 128 {
			return 91 // Bright red
		}
		return 31 // Red
	} else if g > r && g > b {
		if brightness > 128 {
			return 92 // Bright green
		}
		return 32 // Green
	} else if b > r && b > g {
		if brightness > 128 {
			return 94 // Bright blue
		}
		return 34 // Blue
	} else if r > 128 && g > 128 && b < 100 {
		if brightness > 180 {
			return 93 // Bright yellow
		}
		return 33 // Yellow
	} else if r > 128 && b > 128 && g < 100 {
		if brightness > 180 {
			return 95 // Bright magenta
		}
		return 35 // Magenta
	} else if g > 128 && b > 128 && r < 100 {
		if brightness > 180 {
			return 96 // Bright cyan
		}
		return 36 // Cyan
	} else if brightness < 64 {
		return 30 // Black
	} else if brightness > 192 {
		return 97 // Bright white
	} else if brightness > 128 {
		return 37 // White
	}
	return 90 // Bright black (gray)
}

// rgbToANSI256 converts RGB to the closest 256-color palette index
func rgbToANSI256(r, g, b int) int {
	// Check if it's a gray
	if abs(r-g) < 10 && abs(g-b) < 10 && abs(b-r) < 10 {
		// Use grayscale ramp (232-255)
		if r < 8 {
			return 16 // Black
		}
		if r > 247 {
			return 231 // White
		}
		return 232 + (r-8)/10
	}

	// Use 6x6x6 color cube (16-231)
	r6 := (r * 6) / 256
	g6 := (g * 6) / 256
	b6 := (b * 6) / 256

	return 16 + 36*r6 + 6*g6 + b6
}

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
