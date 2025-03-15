package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

var mapNext string
var mapPrevious any

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

type Location struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous any    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func commandMap() error {
	url := ""
	if mapNext != "" {
		url = mapNext
	} else {
		url = "https://pokeapi.co/api/v2/location-area/"
	}
	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Get error: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Reader error")
		return fmt.Errorf("Reader error: %v", err)
	}

	var locations Location
	if err := json.Unmarshal(body, &locations); err != nil {
		fmt.Printf("locations error: %v", err)
		return fmt.Errorf("Unmarshal error: %v", err)
	}
	mapNext = locations.Next
	mapPrevious = locations.Previous
	for i := 0; i < len(locations.Results); i++ {
		fmt.Printf("%s\n", locations.Results[i].Name)
	}

	return nil
}

func commandMapB() error {
	url := ""
	if mapPrevious != nil {
		url = mapPrevious.(string)
	} else {
		fmt.Printf("end of map\n")
		return nil
	}
	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Get error: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Reader error")
		return fmt.Errorf("Reader error: %v", err)
	}

	var locations Location
	if err := json.Unmarshal(body, &locations); err != nil {
		fmt.Printf("locations error: %v", err)
		return fmt.Errorf("Unmarshal error: %v", err)
	}
	mapNext = locations.Next
	mapPrevious = locations.Previous
	for i := 0; i < len(locations.Results); i++ {
		fmt.Printf("%s\n", locations.Results[i].Name)
	}
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
