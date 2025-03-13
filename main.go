package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for true {
		fmt.Print("Pokedex >")
		scanner.Scan()
		command := scanner.Text()
		clean := cleanInput(command)
		fmt.Printf("Your command was: %s\n", clean[0])
	}
}

func cleanInput(text string) []string {
	textStrip := strings.TrimSpace(text)
	textStrip = strings.ToLower(textStrip)
	result := strings.Split(textStrip, " ")
	return result
}
