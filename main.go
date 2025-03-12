package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Printf("Hello, World!")
}

func cleanInput(text string) []string {
	textStrip := strings.TrimSpace(text)
	textStrip = strings.ToLower(textStrip)
	result := strings.Split(textStrip, " ")
	return result
}
