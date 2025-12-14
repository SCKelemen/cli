package renderer

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/term"
)

// TerminalDimensions represents the complete dimensions of the terminal
type TerminalDimensions struct {
	Columns int // Character columns (e.g., 80)
	Rows    int // Character rows (e.g., 24)

	PixelWidth  int // Total width in pixels (0 if unsupported)
	PixelHeight int // Total height in pixels (0 if unsupported)

	CellWidth  float64 // Width of one character in pixels
	CellHeight float64 // Height of one character in pixels

	HasPixelSupport bool // Whether terminal supports pixel queries
}

// QueryTerminalDimensions gets comprehensive terminal dimensions including pixel sizes
func QueryTerminalDimensions(columns, rows int) TerminalDimensions {
	dims := TerminalDimensions{
		Columns: columns,
		Rows:    rows,
	}

	// Try to query pixel dimensions
	pixelWidth, pixelHeight, cellWidth, cellHeight, ok := queryPixelDimensions()
	if ok {
		dims.PixelWidth = pixelWidth
		dims.PixelHeight = pixelHeight
		dims.CellWidth = cellWidth
		dims.CellHeight = cellHeight
		dims.HasPixelSupport = true
	} else {
		// Fallback: extrapolate from typical terminal fonts
		dims.CellWidth = 9.0  // Typical monospace width (8-10px)
		dims.CellHeight = 18.0 // Typical line height (16-20px)
		dims.PixelWidth = int(dims.CellWidth * float64(columns))
		dims.PixelHeight = int(dims.CellHeight * float64(rows))
		dims.HasPixelSupport = false
	}

	return dims
}

// queryPixelDimensions queries the terminal for pixel dimensions using CSI sequences
// Returns: pixelWidth, pixelHeight, cellWidth, cellHeight, success
func queryPixelDimensions() (int, int, float64, float64, bool) {
	// Check if we're in a terminal
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		return 0, 0, 0, 0, false
	}

	// Save current terminal state
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return 0, 0, 0, 0, false
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	// Query window pixel size: CSI 14 t
	// Terminal should respond with: CSI 4 ; height ; width t
	fmt.Fprint(os.Stdout, "\x1b[14t")

	// Read response with timeout
	response := readCSIResponse(100 * time.Millisecond)
	if response == "" {
		return 0, 0, 0, 0, false
	}

	// Parse response: ESC [ 4 ; height ; width t
	pixelWidth, pixelHeight, ok := parsePixelSizeResponse(response)
	if !ok {
		return 0, 0, 0, 0, false
	}

	// Query character size: CSI 18 t
	// Terminal should respond with: CSI 8 ; rows ; columns t
	fmt.Fprint(os.Stdout, "\x1b[18t")

	response = readCSIResponse(100 * time.Millisecond)
	if response == "" {
		// We have pixel size but not character count, calculate from standard dimensions
		return pixelWidth, pixelHeight, 9.0, 18.0, true
	}

	columns, rows, ok := parseCharSizeResponse(response)
	if !ok || columns == 0 || rows == 0 {
		return pixelWidth, pixelHeight, 9.0, 18.0, true
	}

	// Calculate cell dimensions
	cellWidth := float64(pixelWidth) / float64(columns)
	cellHeight := float64(pixelHeight) / float64(rows)

	return pixelWidth, pixelHeight, cellWidth, cellHeight, true
}

// readCSIResponse reads a CSI response from stdin with timeout
func readCSIResponse(timeout time.Duration) string {
	result := make(chan string, 1)

	go func() {
		buf := make([]byte, 32)
		var response strings.Builder

		// Read until we get 't' (end of response) or timeout
		for {
			n, err := os.Stdin.Read(buf)
			if err != nil || n == 0 {
				break
			}

			response.Write(buf[:n])

			// Check if we've received the complete response (ends with 't')
			if buf[n-1] == 't' {
				break
			}

			// Safety: don't read more than 32 bytes
			if response.Len() >= 32 {
				break
			}
		}

		result <- response.String()
	}()

	select {
	case resp := <-result:
		return resp
	case <-time.After(timeout):
		return ""
	}
}

// parsePixelSizeResponse parses CSI 4 ; height ; width t response
func parsePixelSizeResponse(response string) (width, height int, ok bool) {
	// Expected format: ESC [ 4 ; height ; width t
	// Strip ESC [ prefix
	response = strings.TrimPrefix(response, "\x1b[")
	response = strings.TrimPrefix(response, "\033[")

	// Strip 't' suffix
	response = strings.TrimSuffix(response, "t")

	parts := strings.Split(response, ";")
	if len(parts) != 3 {
		return 0, 0, false
	}

	// First part should be "4"
	if parts[0] != "4" {
		return 0, 0, false
	}

	height, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return 0, 0, false
	}

	width, err = strconv.Atoi(strings.TrimSpace(parts[2]))
	if err != nil {
		return 0, 0, false
	}

	return width, height, true
}

// parseCharSizeResponse parses CSI 8 ; rows ; columns t response
func parseCharSizeResponse(response string) (columns, rows int, ok bool) {
	// Expected format: ESC [ 8 ; rows ; columns t
	// Strip ESC [ prefix
	response = strings.TrimPrefix(response, "\x1b[")
	response = strings.TrimPrefix(response, "\033[")

	// Strip 't' suffix
	response = strings.TrimSuffix(response, "t")

	parts := strings.Split(response, ";")
	if len(parts) != 3 {
		return 0, 0, false
	}

	// First part should be "8"
	if parts[0] != "8" {
		return 0, 0, false
	}

	rows, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return 0, 0, false
	}

	columns, err = strconv.Atoi(strings.TrimSpace(parts[2]))
	if err != nil {
		return 0, 0, false
	}

	return columns, rows, true
}

// String returns a human-readable representation of dimensions
func (d TerminalDimensions) String() string {
	if d.HasPixelSupport {
		return fmt.Sprintf("%dx%d chars (%dx%d pixels, %.1fx%.1f per cell)",
			d.Columns, d.Rows, d.PixelWidth, d.PixelHeight, d.CellWidth, d.CellHeight)
	}
	return fmt.Sprintf("%dx%d chars (%.1fx%.1f per cell, estimated)",
		d.Columns, d.Rows, d.CellWidth, d.CellHeight)
}
