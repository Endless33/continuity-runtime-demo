package main

import "fmt"

func main() {
	fmt.Println("CONTINUITY RUNTIME DEMO")
	fmt.Println()

	fmt.Println("[EVENT] WiFi failed")
	fmt.Println("[DECISION] migrate=true (margin=87.8, confidence=1.00, reason=better_path)")
	fmt.Println("[AUTHORITY] epoch 2 granted to 5G")
	fmt.Println("[CHECK] stale WiFi rejected")
	fmt.Println("[RESULT] session continues")
}