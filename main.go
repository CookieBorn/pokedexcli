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
		calledCommand := false
		for com, cli := range callbackMap {
			if com == clean[0] {
				calledCommand = true
				cli.callback()
			}
		}
		if calledCommand == false {
			fmt.Print("Unknown command\n")
		}

	}
}

func cleanInput(text string) []string {
	textStrip := strings.TrimSpace(text)
	textStrip = strings.ToLower(textStrip)
	result := strings.Split(textStrip, " ")
	return result
}

func commandExit() error {
	fmt.Printf("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
	for _, cli := range callbackMap {
		fmt.Printf("%s: %s\n", cli.name, cli.description)
	}
	return nil
}

func commandMap() error {
	return nil
}

func commandMapB() error {
	return nil
}

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var callbackMap map[string]cliCommand

func init() {
	callbackMap = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "show the next 20 locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "show the previous 20 locations",
			callback:    commandMapB,
		},
	}
}
