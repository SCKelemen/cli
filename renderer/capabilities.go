package renderer

import (
	"os"
	"strings"
)

// ColorMode represents the terminal's color capabilities
type ColorMode int

const (
	ColorModeNone   ColorMode = iota // No color support
	ColorMode16                      // 16 ANSI colors
	ColorMode256                     // 256 colors
	ColorModeTrueColor               // 24-bit RGB (16.7M colors)
)

// TerminalCapabilities holds information about terminal capabilities
type TerminalCapabilities struct {
	ColorMode   ColorMode
	IsTTY       bool
	SupportsAlt bool // Alternate screen buffer
}

// DetectCapabilities detects the terminal's capabilities
func DetectCapabilities() *TerminalCapabilities {
	caps := &TerminalCapabilities{
		ColorMode:   ColorModeNone,
		IsTTY:       false,
		SupportsAlt: true,
	}

	// Check if stdout is a terminal
	if fileInfo, err := os.Stdout.Stat(); err == nil {
		caps.IsTTY = (fileInfo.Mode() & os.ModeCharDevice) != 0
	}

	if !caps.IsTTY {
		return caps
	}

	// Detect color mode
	caps.ColorMode = detectColorMode()

	return caps
}

func detectColorMode() ColorMode {
	// Check TERM_PROGRAM for specific terminal applications
	termProgram := os.Getenv("TERM_PROGRAM")

	// Apple's Terminal.app claims true color support but only supports 16 colors
	if termProgram == "Apple_Terminal" {
		return ColorMode16
	}

	// Check COLORTERM environment variable (most reliable for true color)
	colorTerm := os.Getenv("COLORTERM")
	if colorTerm == "truecolor" || colorTerm == "24bit" {
		return ColorModeTrueColor
	}

	// Check TERM environment variable
	term := os.Getenv("TERM")

	// True color terminals
	if strings.Contains(term, "truecolor") || strings.Contains(term, "24bit") {
		return ColorModeTrueColor
	}

	// iTerm2, Alacritty, Kitty support true color
	if strings.Contains(term, "iterm") ||
		strings.Contains(term, "alacritty") ||
		strings.Contains(term, "kitty") {
		return ColorModeTrueColor
	}

	// 256 color terminals
	if strings.Contains(term, "256color") {
		return ColorMode256
	}

	// xterm variants usually support 256 colors
	if strings.Contains(term, "xterm") {
		return ColorMode256
	}

	// Basic color support
	if term != "" && term != "dumb" {
		return ColorMode16
	}

	return ColorModeNone
}

// String returns a human-readable description of the color mode
func (cm ColorMode) String() string {
	switch cm {
	case ColorModeNone:
		return "No color"
	case ColorMode16:
		return "16 colors"
	case ColorMode256:
		return "256 colors"
	case ColorModeTrueColor:
		return "True color (24-bit)"
	default:
		return "Unknown"
	}
}
