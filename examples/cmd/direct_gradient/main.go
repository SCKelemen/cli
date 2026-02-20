package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/SCKelemen/color"
)

func main() {
	// Clear screen
	fmt.Print("\033[2J\033[H")

	fmt.Print("=== Direct ANSI Gradient Test ===\n")

	// Test 1: Simple gradient using direct ANSI codes
	fmt.Println("1. Blue to Red (Direct ANSI):")
	for i := 0; i < 70; i++ {
		t := float64(i) / 69.0
		hue := 240 + (t * 150)
		if hue > 360 {
			hue -= 360
		}
		colorStr := fmt.Sprintf("oklch(0.6 0.2 %.0f)", hue)
		c, _ := color.ParseColor(colorStr)
		r, g, b, _ := c.RGBA()
		fmt.Printf("\033[48;2;%d;%d;%dm ", int(r*255), int(g*255), int(b*255))
	}
	fmt.Print("\033[0m\n")

	fmt.Println()

	// Test 2: Rainbow
	fmt.Println("2. Rainbow (Hue 0° → 360°):")
	for i := 0; i < 70; i++ {
		t := float64(i) / 69.0
		hue := t * 360
		colorStr := fmt.Sprintf("oklch(0.65 0.2 %.0f)", hue)
		c, _ := color.ParseColor(colorStr)
		r, g, b, _ := c.RGBA()
		fmt.Printf("\033[48;2;%d;%d;%dm ", int(r*255), int(g*255), int(b*255))
	}
	fmt.Print("\033[0m\n")

	fmt.Println()

	// Test 3: Using █ blocks for double height
	fmt.Println("3. Gradient with blocks:")
	for i := 0; i < 70; i++ {
		t := float64(i) / 69.0
		hue := 240 + (t * 150)
		if hue > 360 {
			hue -= 360
		}
		colorStr := fmt.Sprintf("oklch(0.6 0.2 %.0f)", hue)
		c, _ := color.ParseColor(colorStr)
		r, g, b, _ := c.RGBA()
		fmt.Printf("\033[48;2;%d;%d;%dm \033[0m", int(r*255), int(g*255), int(b*255))
	}
	fmt.Println()
	for i := 0; i < 70; i++ {
		t := float64(i) / 69.0
		hue := 240 + (t * 150)
		if hue > 360 {
			hue -= 360
		}
		colorStr := fmt.Sprintf("oklch(0.6 0.2 %.0f)", hue)
		c, _ := color.ParseColor(colorStr)
		r, g, b, _ := c.RGBA()
		fmt.Printf("\033[48;2;%d;%d;%dm \033[0m", int(r*255), int(g*255), int(b*255))
	}
	fmt.Println()

	fmt.Print("\n=== Terminal Info ===")
	fmt.Printf("TERM=%s\n", os.Getenv("TERM"))
	fmt.Printf("COLORTERM=%s\n", os.Getenv("COLORTERM"))

	fmt.Print("\n=== Expected Result ===")
	fmt.Println("✓ Line 1: Should show smooth blue → purple → red gradient")
	fmt.Println("✓ Line 2: Should show full rainbow spectrum")
	fmt.Println("✓ Line 3: Should show gradient in two rows")

	fmt.Print("\nIf you see solid colors (not gradients), your terminal")
	fmt.Println("may not support 24-bit true color.")

	// Wait for user
	fmt.Print("\nPress Enter to exit...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
