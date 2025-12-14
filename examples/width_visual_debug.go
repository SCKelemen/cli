package main

import (
	"fmt"
)

func main() {
	fmt.Println("=== Visual Width Test ===")
	fmt.Println("Each line should end with | at column 10:\n")

	fmt.Println("12345678|")
	fmt.Println("• test  |")
	fmt.Println("✓ test  |")
	fmt.Println("a test  |")

	fmt.Println("\nIf • and ✓ lines don't align, they're taking 2 columns")
	fmt.Println("\nDetailed test:")
	fmt.Println("[•]12345|")
	fmt.Println("[✓]12345|")
	fmt.Println("[a]12345|")
}
