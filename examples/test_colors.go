package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("=== Terminal Color Support Test ===\n")

	// Check if stdout is a terminal
	if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) == 0 {
		fmt.Println("⚠️  Output is not a TTY (terminal)")
		fmt.Println("   Try running directly in your terminal, not piped/redirected")
		fmt.Println()
	}

	// Test basic ANSI colors (16 colors)
	fmt.Println("1. Basic ANSI Colors (16 colors):")
	for i := 30; i <= 37; i++ {
		fmt.Printf("\033[%dm  Color %d  \033[0m ", i, i)
	}
	fmt.Println()

	// Test 256 colors
	fmt.Println("\n2. 256 Color Palette:")
	for i := 0; i < 16; i++ {
		fmt.Printf("\033[48;5;%dm  \033[0m", i)
	}
	fmt.Println()

	// Test true color (24-bit RGB)
	fmt.Println("\n3. True Color (24-bit RGB):")
	fmt.Print("Should see: ")
	fmt.Print("\033[48;2;255;0;0m   RED   \033[0m ")
	fmt.Print("\033[48;2;0;255;0m  GREEN  \033[0m ")
	fmt.Print("\033[48;2;0;0;255m   BLUE  \033[0m")
	fmt.Println()

	// Test gradient
	fmt.Println("\n4. Gradient Test:")
	for i := 0; i < 50; i++ {
		r := int(float64(i) / 50.0 * 255)
		fmt.Printf("\033[48;2;%d;0;%dm \033[0m", r, 255-r)
	}
	fmt.Println()

	fmt.Println("\n=== Diagnostic Info ===")
	fmt.Printf("TERM=%s\n", os.Getenv("TERM"))
	fmt.Printf("COLORTERM=%s\n", os.Getenv("COLORTERM"))

	fmt.Println("\n=== What Should You See? ===")
	fmt.Println("✓ Section 1: Different colored backgrounds")
	fmt.Println("✓ Section 2: A row of colored squares")
	fmt.Println("✓ Section 3: Red, green, and blue blocks")
	fmt.Println("✓ Section 4: A smooth gradient from red to blue")

	fmt.Println("\n=== If You Don't See Colors ===")
	fmt.Println("Your terminal may not support true color. Try:")
	fmt.Println("• iTerm2 (macOS) - Full true color support")
	fmt.Println("• Alacritty - Modern, fast, true color")
	fmt.Println("• Windows Terminal - True color support")
	fmt.Println("• Kitty - Advanced terminal emulator")
	fmt.Println("• Set COLORTERM=truecolor in your environment")
}
