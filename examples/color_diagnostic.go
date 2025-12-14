package main

import (
	"fmt"
	"os"

	"github.com/SCKelemen/cli/renderer"
)

func main() {
	fmt.Println("=== Terminal Color Detection ===")
	fmt.Printf("TERM=%s\n", os.Getenv("TERM"))
	fmt.Printf("COLORTERM=%s\n", os.Getenv("COLORTERM"))

	caps := renderer.DetectCapabilities()
	fmt.Printf("Detected mode: %s\n", caps.ColorMode)
	fmt.Printf("Is TTY: %t\n\n", caps.IsTTY)

	// Test each color mode
	modes := []renderer.ColorMode{
		renderer.ColorMode16,
		renderer.ColorMode256,
		renderer.ColorModeTrueColor,
	}

	for _, mode := range modes {
		fmt.Printf("%s test:\n", mode)
		r := renderer.NewANSIRendererWithMode(mode)

		// Test red background
		fmt.Print(r.Reset())
		fmt.Print("\x1b[41m")
		fmt.Print("  RED  ")
		fmt.Print(r.Reset())
		fmt.Print(" ")

		// Test green background
		fmt.Print("\x1b[42m")
		fmt.Print(" GREEN ")
		fmt.Print(r.Reset())
		fmt.Print(" ")

		// Test blue background
		fmt.Print("\x1b[44m")
		fmt.Print(" BLUE  ")
		fmt.Print(r.Reset())
		fmt.Println()
	}
}
