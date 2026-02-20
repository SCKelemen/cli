package main

import (
	"fmt"

	"github.com/mattn/go-runewidth"
)

func main() {
	fmt.Println("=== Box Drawing Character Widths ===\n")

	chars := []rune{
		'╭', '─', '╮',
		'│', ' ', '│',
		'╰', '─', '╯',
		'•', '✓',
	}

	for _, ch := range chars {
		w := runewidth.RuneWidth(ch)
		fmt.Printf("'%c' (U+%04X) = width %d\n", ch, ch, w)
	}

	fmt.Println("\n=== Test Rendering ===")
	fmt.Println("╭──────╮")
	fmt.Println("│ Text │")
	fmt.Println("╰──────╯")
}
