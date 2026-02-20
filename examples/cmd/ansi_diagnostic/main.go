package main

import (
	"fmt"
)

func main() {
	fmt.Println("=== Direct ANSI Color Test ===\n")

	// Test 16 basic colors (backgrounds)
	fmt.Println("16-color ANSI backgrounds:")
	colors := []struct {
		name string
		code int
	}{
		{"Black", 40},
		{"Red", 41},
		{"Green", 42},
		{"Yellow", 43},
		{"Blue", 44},
		{"Magenta", 45},
		{"Cyan", 46},
		{"White", 47},
		{"Bright Black", 100},
		{"Bright Red", 101},
		{"Bright Green", 102},
		{"Bright Yellow", 103},
		{"Bright Blue", 104},
		{"Bright Magenta", 105},
		{"Bright Cyan", 106},
		{"Bright White", 107},
	}

	for _, c := range colors {
		fmt.Printf("\x1b[%dm %-15s \x1b[0m ", c.code, c.name)
		if (c.code-40+1)%4 == 0 {
			fmt.Println()
		}
	}

	fmt.Println("\n\n256-color test (should show gradient):")
	for i := 16; i < 52; i++ {
		fmt.Printf("\x1b[48;5;%dm  \x1b[0m", i)
	}
	fmt.Println()

	fmt.Println("\n\nTruecolor test (should show smooth gradient):")
	for i := 0; i < 50; i++ {
		r := (i * 255) / 50
		fmt.Printf("\x1b[48;2;%d;0;%dm  \x1b[0m", 255-r, r)
	}
	fmt.Println()
	fmt.Println()
}
